package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/tikalk/resource-manager/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/types"
)

// ObjectHandler manage a single object like deployment, namespace, etc...
// according to the action definition provided by user like "delete" / "patch" an object
type ObjectHandler struct {
	resourceManager *v1alpha1.ResourceManager
	object          interface{}
	fullname        types.NamespacedName
	creationTime    time.Time
	stopper         chan struct{}
	clientset       *kubernetes.Clientset
	log             logr.Logger
}

// NewObjectHandler create a new ObjectHandler to manage a single kubernetes object
func NewObjectHandler(resourceManager *v1alpha1.ResourceManager, obj interface{}, clientset *kubernetes.Clientset, log logr.Logger) (*ObjectHandler, error) {
	// extract the NamespacedName of the object for storage
	fullName, err := extractFullname(resourceManager.Spec.ResourceKind, obj)
	if err != nil {
		return nil, err
	}

	creationTime, err := extractCreationTime(resourceManager.Spec.ResourceKind, obj)
	if err != nil {
		return nil, err
	}

	// return the object handler
	objectHandler := &ObjectHandler{
		object:          obj,
		fullname:        fullName,
		creationTime:    creationTime,
		stopper:         make(chan struct{}),
		resourceManager: resourceManager,
		clientset:       clientset,
		log:             log,
	}
	return objectHandler, nil
}

// extractFullname extract the full name of the object according to object kind
func extractFullname(kind string, obj interface{}) (fullname types.NamespacedName, err error) {
	switch kind {
	case "Namespace":
		fullname = types.NamespacedName{Name: obj.(*v1.Namespace).Name, Namespace: obj.(*v1.Namespace).Namespace}
	case "Deployment":
		fullname = types.NamespacedName{Name: obj.(*appsv1.Deployment).Name, Namespace: obj.(*appsv1.Deployment).Namespace}
	default:
		err = fmt.Errorf("extractFullname error: unxpected object kind <%s>", kind)
	}
	return fullname, err
}

// extractCreationTime extract the creation time of the object according to object kind
func extractCreationTime(kind string, obj interface{}) (time time.Time, err error) {
	switch kind {
	case "Namespace":
		time = obj.(*v1.Namespace).ObjectMeta.CreationTimestamp.Time
	case "Deployment":
		time = obj.(*appsv1.Deployment).ObjectMeta.CreationTimestamp.Time
	default:
		err = fmt.Errorf("extractCreationTime: unxpected object kind <%s>", kind)
	}
	return time, err
}

// performObjectAction executes the desired action on an object
func (h *ObjectHandler) performObjectAction() (err error) {
	switch h.resourceManager.Spec.Action {
	case "delete":
		err = h.performObjectDelete()
		break
	case "patch":
		err = h.performObjectPatch()
		break
	default:
		err = errors.New(fmt.Sprintf("objectAction: unexpected action %s", h.resourceManager.Spec.Action))
	}
	return err
}

// performObjectDelete delete a single object
func (h *ObjectHandler) performObjectDelete() (err error) {
	var opts metav1.DeleteOptions
	switch h.resourceManager.Spec.ResourceKind {
	case "Namespace":
		err = h.clientset.CoreV1().Namespaces().Delete(context.Background(), h.fullname.Name, opts)
	case "Deployment":
		err = h.clientset.AppsV1().Deployments(h.fullname.Namespace).Delete(context.Background(), h.fullname.Name, opts)
	default:
		err = fmt.Errorf("objectDelete: unxpected object kind <%s>", h.resourceManager.Spec.ResourceKind)
	}
	return err
}

type patchUInt32Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

// performObjectPatch patch a single object
func (h *ObjectHandler) performObjectPatch() (err error) {
	//var pt types.PatchType

	//data := fmt.Sprintf(`{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}`, time.Now().String())
	var data string
	data = h.resourceManager.Spec.ActionParam

	switch h.resourceManager.Spec.ResourceKind {
	case "Namespace":
		_, err = h.clientset.CoreV1().Namespaces().Patch(context.Background(), h.fullname.Name, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	case "Deployment":
		_, err = h.clientset.AppsV1().Deployments(h.fullname.Namespace).Patch(context.Background(), h.fullname.Name, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	default:
		err = fmt.Errorf("objectDelete: unxpected object kind <%s>", h.resourceManager.Spec.ResourceKind)
	}
	return err
}

// Run calculates the expiration time of an object and perform the desired action when the time arrives
func (h *ObjectHandler) Run() {
	var wait time.Duration

	cond := h.resourceManager.Spec.Condition
	if cond.ExpireAfter != "" {
		expireAfter, err := time.ParseDuration(cond.ExpireAfter)
		if err != nil {
			h.log.Error(err, trace(fmt.Sprintf("cannot parse ExpireAfter parameter <%s>. object handler <%s> aborted", cond.ExpireAfter, h.fullname)))
			return
		}
		age := time.Now().Sub(h.creationTime)
		wait = expireAfter - age

		h.log.Info(trace(fmt.Sprintf("object age expiration <%s> after <%s> age <%s> wait <%s>", h.fullname, expireAfter.String(), age.String(), wait.String())))
	} else if cond.ExpireAt != "" {
		expireAt, err := time.Parse("15:04", cond.ExpireAt)
		if err != nil {
			h.log.Error(err, trace(fmt.Sprintf("Failed to parse <%s>. object handler <%s> aborted.", cond.ExpireAt, h.fullname)))
			return
		}

		now := time.Now()
		if expireAt.Hour()*60+expireAt.Minute() > now.Hour()*60+now.Minute() {
			// Today
			wait = time.Date(now.Year(), now.Month(), now.Day(), expireAt.Hour(), expireAt.Minute(), 0, 0, now.Location()).Sub(now)
		} else {
			// Tomorrow
			tomorrow := now.Add(24 * time.Hour)
			wait = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), expireAt.Hour(), expireAt.Minute(), 0, 0, tomorrow.Location()).Sub(now)
		}

		h.log.Info(trace(fmt.Sprintf("object time expiration <%s> expireAt <%s> now <%s> wait <%s>", h.fullname, expireAt.String(), now.String(), wait)))
	} else {
		h.log.Error(errors.New("expiration is not configured"), trace(fmt.Sprintf("object handler <%s> aborted", h.fullname)))
		return
	}

	if wait <= 0 {
		h.log.Info(trace(fmt.Sprintf("object already expired <%s>", h.fullname)))
	} else {
		select {
		case <-h.stopper:
			h.log.Info(trace(fmt.Sprintf("h aborted for object<%s>", h.fullname)))
			return
		case <-time.After(wait):
			h.log.Info(trace(fmt.Sprintf("object expired <%s>", h.fullname)))
			break
		}
	}

	if h.resourceManager.Spec.DryRun {
		h.log.Info(trace(fmt.Sprintf("dry-run performing object <%s> action <%s> ", h.fullname, h.resourceManager.Spec.Action)))
	} else {
		h.log.Info(trace(fmt.Sprintf("performing object <%s> action <%s>...", h.fullname, h.resourceManager.Spec.Action)))
		err := h.performObjectAction()
		if err != nil {
			h.log.Error(err, trace(fmt.Sprintf("object <%s> action <%s> failed", h.fullname, h.resourceManager.Spec.Action)))
		} else {
			h.log.Info(trace(fmt.Sprintf("object <%s> action <%s> finished", h.fullname, h.resourceManager.Spec.Action)))
		}

	}

}

// Stop will be called, When the ObjectHandler requires to stop.
func (h *ObjectHandler) Stop() {
	close(h.stopper)
}
