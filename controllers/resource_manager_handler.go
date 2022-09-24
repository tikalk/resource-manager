package controllers

import (
	"fmt"
	"github.com/go-logr/logr"
	v1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type ResourceManagerHandler struct {
	resourceManager *v1alpha1.ResourceManager
	informer        cache.SharedIndexInformer
	//factory         informers.SharedInformerFactory
	objHandlers map[types.NamespacedName]*ObjectHandler

	stopper   chan struct{}
	clientset *kubernetes.Clientset
	log       logr.Logger
}

func NewResourceManagerHandler(resourceManager *v1alpha1.ResourceManager, clientset *kubernetes.Clientset, log logr.Logger) (*ResourceManagerHandler, error) {
	selector, _ := metav1.LabelSelectorAsSelector(resourceManager.Spec.Selector)
	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(resourceManager.Namespace), labelOptions)

	informer, err := createInformer(factory, resourceManager.Spec.ResourceKind)
	if err != nil {
		return nil, err
	}

	return &ResourceManagerHandler{
		resourceManager: resourceManager,
		informer:        informer,
		//factory:         factory,
		objHandlers: make(map[types.NamespacedName]*ObjectHandler),
		stopper:     make(chan struct{}),
		clientset:   clientset,
		log:         log,
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

func (resourceManagerHandler *ResourceManagerHandler) addObjHandelr(objHandler *ObjectHandler) {
	resourceManagerHandler.objHandlers[objHandler.fullname] = objHandler
}

func (resourceManagerHandler *ResourceManagerHandler) removeObjHandelr(fullname types.NamespacedName) {
	if _, ok := resourceManagerHandler.objHandlers[fullname]; ok {
		resourceManagerHandler.objHandlers[fullname].Stop()
		delete(resourceManagerHandler.objHandlers, fullname)
	}
}

func (resourceManagerHandler *ResourceManagerHandler) Run() error {

	// TODO: listen to events
	resourceManagerHandler.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			objHandler, err := NewObjectHandler(resourceManagerHandler.resourceManager, obj, resourceManagerHandler.clientset, resourceManagerHandler.log)
			if err != nil {
				resourceManagerHandler.log.Error(err, fmt.Sprintf("NewObjectHandler handler creating failed with error <%s>.", err))
				return
			}
			if objHandler.terminating {
				resourceManagerHandler.log.Info(trace(fmt.Sprintf("Object adding ignored: <%s> Terminating <%b>", objHandler.fullname, objHandler.terminating)))
				return
			}
			resourceManagerHandler.log.Info(trace(fmt.Sprintf("Adding object handler: <%s>", objHandler.fullname)))
			resourceManagerHandler.addObjHandelr(objHandler)
			go objHandler.Run()
			// TODO: handle terminating state  - obj.(*v1.Namespace).Status.Phase == "Terminating"
		},
		UpdateFunc: func(oldObj interface{}, obj interface{}) {
			objHandler, err := NewObjectHandler(resourceManagerHandler.resourceManager, obj, resourceManagerHandler.clientset, resourceManagerHandler.log)
			if err != nil {
				resourceManagerHandler.log.Error(err, fmt.Sprintf("NewObjectHandler handler creating failed with error <%s>.", err))
				return
			}
			if objHandler.terminating {
				resourceManagerHandler.log.Info(trace(fmt.Sprintf("Object recreating ignored: <%s> Terminating <%b>", objHandler.fullname, objHandler.terminating)))
				return
			}
			resourceManagerHandler.log.Info(trace(fmt.Sprintf("Recreating object handler: <%s>", objHandler.fullname)))
			resourceManagerHandler.removeObjHandelr(objHandler.fullname)
			resourceManagerHandler.addObjHandelr(objHandler)
			go objHandler.Run()
		},
		DeleteFunc: func(obj interface{}) {
			objHandler, err := NewObjectHandler(resourceManagerHandler.resourceManager, obj, resourceManagerHandler.clientset, resourceManagerHandler.log)
			if err != nil {
				resourceManagerHandler.log.Error(err, fmt.Sprintf("NewObjectHandler handler creating failed with error <%s>.", err))
				return
			}
			resourceManagerHandler.log.Info(trace(fmt.Sprintf("Deleting object handler: <%s>", objHandler.fullname)))
			resourceManagerHandler.removeObjHandelr(objHandler.fullname)
		},
	})
	// start the informer
	go resourceManagerHandler.informer.Run(resourceManagerHandler.stopper)

	return nil
}

func (resourceManagerHandler *ResourceManagerHandler) Stop() {
	// stop channel
	close(resourceManagerHandler.stopper)

	// TODO: cleanup
	for _, h := range resourceManagerHandler.objHandlers {
		h.Stop()
	}
}
