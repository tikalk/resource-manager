package controllers

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/tikalk/resource-manager/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Represents object namespace with managed objects
type ObjectNamespaceHandler struct {
	resourceManager *v1alpha1.ResourceManager
	namespaceName   string
	objectsInformer cache.SharedIndexInformer
	objHandlers     map[types.NamespacedName]*ObjectHandler
	stopper         chan struct{}
	clientset       *kubernetes.Clientset
	log             logr.Logger
}

func NewObjectNamespaceHandler(resourceManager *v1alpha1.ResourceManager, namespaceName string, clientset *kubernetes.Clientset, log logr.Logger) (*ObjectNamespaceHandler, error) {
	selector, _ := metav1.LabelSelectorAsSelector(resourceManager.Spec.Selector)
	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(namespaceName), labelOptions)

	objectsInformer, err := createInformer(factory, resourceManager.Spec.ResourceKind)
	if err != nil {
		return nil, err
	}

	return &ObjectNamespaceHandler{
		resourceManager: resourceManager,
		namespaceName:   namespaceName,
		objectsInformer: objectsInformer,
		objHandlers:     make(map[types.NamespacedName]*ObjectHandler),
		stopper:         make(chan struct{}),
		clientset:       clientset,
		log:             log,
	}, nil
}

func createInformer(factory informers.SharedInformerFactory, kind string) (informer cache.SharedIndexInformer, err error) {
	switch kind {
	case "Deployment":
		informer = factory.Apps().V1().Deployments().Informer()
	case "Namespace":
		informer = factory.Core().V1().Namespaces().Informer()
	default:
		err = fmt.Errorf("invalid kind %s when getting an informer.", kind)
	}
	return informer, err
}

func (objectNamespaceHandler *ObjectNamespaceHandler) addObjHandler(objHandler *ObjectHandler) {
	if _, ok := objectNamespaceHandler.objHandlers[objHandler.fullname]; ok {
		objectNamespaceHandler.log.Error(errors.New("addObjHandler failed"), trace(fmt.Sprintf("object handler already registered <%s>.", objHandler.fullname)))
		return
	}

	objectNamespaceHandler.objHandlers[objHandler.fullname] = objHandler
}

func (objectNamespaceHandler *ObjectNamespaceHandler) removeObjHandelr(fullname types.NamespacedName) {
	if _, ok := objectNamespaceHandler.objHandlers[fullname]; !ok {
		objectNamespaceHandler.log.Error(errors.New("removeObjHandelr failed"), trace(fmt.Sprintf("object handler already registered <%s>.", fullname)))
		return
	}
	objectNamespaceHandler.objHandlers[fullname].Stop()
	delete(objectNamespaceHandler.objHandlers, fullname)
}

func (objectNamespaceHandler *ObjectNamespaceHandler) Run() error {

	objectNamespaceHandler.objectsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			objHandler, err := NewObjectHandler(objectNamespaceHandler.resourceManager, obj, objectNamespaceHandler.clientset, objectNamespaceHandler.log)
			if err != nil {
				objectNamespaceHandler.log.Error(err, fmt.Sprintf("NewObjectHandler handler creating failed with error <%s>.", err))
				return
			}
			if objHandler.terminating {
				objectNamespaceHandler.log.Info(trace(fmt.Sprintf("Object adding ignored: <%s> Terminating <%b>", objHandler.fullname, objHandler.terminating)))
				return
			}
			objectNamespaceHandler.log.Info(trace(fmt.Sprintf("Adding object handler: <%s>", objHandler.fullname)))
			objectNamespaceHandler.addObjHandler(objHandler)
			go objHandler.Run()
		},
		//UpdateFunc: func(oldObj interface{}, obj interface{}) {
		//
		//	objHandler, err := NewObjectHandler(objectNamespaceHandler.resourceManager, obj, objectNamespaceHandler.clientset, objectNamespaceHandler.log)
		//	if err != nil {
		//		objectNamespaceHandler.log.Error(err, fmt.Sprintf("NewObjectHandler handler creating failed with error <%s>.", err))
		//		return
		//	}
		//	if objHandler.terminating {
		//		objectNamespaceHandler.log.Info(trace(fmt.Sprintf("Object recreating ignored: <%s> Terminating <%b>", objHandler.fullname, objHandler.terminating)))
		//		return
		//	}
		//	objectNamespaceHandler.log.Info(trace(fmt.Sprintf("Recreating object handler: <%s>", objHandler.fullname)))
		//	objectNamespaceHandler.removeObjHandelr(objHandler.fullname)
		//	objectNamespaceHandler.addObjHandler(objHandler)
		//	go objHandler.Run()
		//},
		DeleteFunc: func(obj interface{}) {
			objHandler, err := NewObjectHandler(objectNamespaceHandler.resourceManager, obj, objectNamespaceHandler.clientset, objectNamespaceHandler.log)
			if err != nil {
				objectNamespaceHandler.log.Error(err, fmt.Sprintf("NewObjectHandler handler creating failed with error <%s>.", err))
				return
			}
			objectNamespaceHandler.log.Info(trace(fmt.Sprintf("Deleting object handler: <%s>", objHandler.fullname)))
			objectNamespaceHandler.removeObjHandelr(objHandler.fullname)
		},
	})
	// start the objectsInformer
	go objectNamespaceHandler.objectsInformer.Run(objectNamespaceHandler.stopper)

	return nil
}

func (objectNamespaceHandler *ObjectNamespaceHandler) Stop() {
	close(objectNamespaceHandler.stopper)
}
