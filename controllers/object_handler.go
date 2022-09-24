package controllers

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/tikalk/resource-manager/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"

	"k8s.io/apimachinery/pkg/types"
)

type ObjectHandler struct {
	resourceManager *v1alpha1.ResourceManager
	object          interface{}
	fullname        types.NamespacedName
	creationTime    time.Time
	terminating     bool
	stopper         chan struct{}
	clientset       *kubernetes.Clientset
	log             logr.Logger
}

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

	terminating, err := extractTerminating(resourceManager.Spec.ResourceKind, obj)
	if err != nil {
		return nil, err
	}

	// return the object handler
	objectHandler := &ObjectHandler{
		object:          obj,
		fullname:        fullName,
		creationTime:    creationTime,
		terminating:     terminating,
		stopper:         make(chan struct{}),
		resourceManager: resourceManager,
		clientset:       clientset,
		log:             log,
	}
	return objectHandler, nil
}

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

func extractTerminating(kind string, obj interface{}) (phase bool, err error) {
	switch kind {
	case "Namespace":
		phase = obj.(*v1.Namespace).Status.Phase == "Terminating"
	case "Deployment":
		phase = obj.(*appsv1.Deployment).Status.String() == "Terminating"
	default:
		err = fmt.Errorf("extractTerminating error: unxpected object kind <%s>", kind)
	}
	return phase, err
}

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

func (objectHandler *ObjectHandler) performObjectAction() (err error) {
	switch objectHandler.resourceManager.Spec.Action {
	case "delete":
		err = objectHandler.performObjectDelete()
		break
	default:
		err = errors.New(fmt.Sprintf("objectAction: unexpected action %s", objectHandler.resourceManager.Spec.Action))
	}
	return err
}

func (objectHandler *ObjectHandler) performObjectDelete() (err error) {
	var opts metav1.DeleteOptions
	switch objectHandler.resourceManager.Spec.ResourceKind {
	case "Namespace":
		err = objectHandler.clientset.CoreV1().Namespaces().Delete(ctx, objectHandler.fullname.Name, opts)
	case "Deployment":
		err = objectHandler.clientset.AppsV1().Deployments(objectHandler.resourceManager.Namespace).Delete(ctx, objectHandler.fullname.Name, opts)
	default:
		err = fmt.Errorf("objectDelete: unxpected object kind <%s>", objectHandler.resourceManager.Spec.ResourceKind)
	}
	return err
}

// func (o *ObjectHandler) getObjType() string {
// 	return o.resourceManager.Spec.ResourceKind
// }

func (objectHandler *ObjectHandler) Run() {
	var wait time.Duration
	if objectHandler.resourceManager.Spec.Condition.Timeframe != "" {
		timeframe, _ := time.ParseDuration(objectHandler.resourceManager.Spec.Condition.Timeframe)
		age := time.Now().Sub(objectHandler.creationTime)
		wait = timeframe - age

		objectHandler.log.Info(trace(fmt.Sprintf("object timeframe expiration <%s> timeframe <%s> age <%s> wait <%s>",
			objectHandler.fullname,
			timeframe.String(),
			age.String(),
			wait.String())))
	} else if objectHandler.resourceManager.Spec.Condition.ExpireAt != "" {
		expireAt, err := time.Parse("15:04", objectHandler.resourceManager.Spec.Condition.ExpireAt)
		if err != nil {
			objectHandler.log.Error(err, fmt.Sprintf("Failed to parse %s", objectHandler.resourceManager.Spec.Condition.ExpireAt))
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

		objectHandler.log.Info(trace(fmt.Sprintf("object time expiration <%s> expireAt <%s> now <%s> wait <%s>",
			objectHandler.fullname,
			expireAt.String(),
			now.String(),
			wait)))
	} else {
		objectHandler.log.Error(errors.New("expiration is not configured"), trace(fmt.Sprintf("object handler <%s> aborted", objectHandler.fullname)))
		return
	}

	if wait <= 0 {
		objectHandler.log.Info(trace(fmt.Sprintf("object already expired <%s>", objectHandler.fullname)))
	} else {
		select {
		case <-objectHandler.stopper:
			objectHandler.log.Info(trace(fmt.Sprintf("objectHandler aborted for object<%s>", objectHandler.fullname)))
			return
		case <-time.After(wait):
			objectHandler.log.Info(trace(fmt.Sprintf("object expired <%s>", objectHandler.fullname)))
			break
		}
	}

	if objectHandler.resourceManager.Spec.DryRun {
		objectHandler.log.Info(trace(fmt.Sprintf("dry-run performing object <%s> action <%s> ", objectHandler.fullname, objectHandler.resourceManager.Spec.Action)))
	} else {
		objectHandler.log.Info(trace(fmt.Sprintf("performing object <%s> action <%s>...", objectHandler.fullname, objectHandler.resourceManager.Spec.Action)))
		err := objectHandler.performObjectAction()
		if err != nil {
			objectHandler.log.Error(err, trace(fmt.Sprintf("object <%s> action <%s> failed", objectHandler.fullname, objectHandler.resourceManager.Spec.Action)))
		} else {
			objectHandler.log.Info(trace(fmt.Sprintf("object <%s> action <%s> finished", objectHandler.fullname, objectHandler.resourceManager.Spec.Action)))
		}

	}
}

func (ObjectHandler *ObjectHandler) Stop() {
	close(ObjectHandler.stopper)
}
