package controllers

import (
	"fmt"
	"github.com/go-logr/logr"
	v1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"reflect"
)

type ResourceManagerHandler struct {
	resourceManager    *v1alpha1.ResourceManager
	namespacesInformer cache.SharedIndexInformer
	namespaceHandlers  map[string]*ObjectNamespaceHandler
	stopper            chan struct{}
	clientset          *kubernetes.Clientset
	log                logr.Logger
}

func NewResourceManagerHandler(resourceManager *v1alpha1.ResourceManager, clientset *kubernetes.Clientset, log logr.Logger) (*ResourceManagerHandler, error) {
	if resourceManager.Spec.NamespaceSelector != nil {
		selector, _ := metav1.LabelSelectorAsSelector(resourceManager.Spec.NamespaceSelector)
		labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
			opts.LabelSelector = selector.String()
		})
		factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, labelOptions)

		namespacesInformer := factory.Core().V1().Namespaces().Informer()

		return &ResourceManagerHandler{
			resourceManager:    resourceManager,
			namespacesInformer: namespacesInformer,
			namespaceHandlers:  make(map[string]*ObjectNamespaceHandler),
			stopper:            make(chan struct{}),
			clientset:          clientset,
			log:                log,
		}, nil
	} else {
		return &ResourceManagerHandler{
			resourceManager:    resourceManager,
			namespacesInformer: nil,
			namespaceHandlers:  make(map[string]*ObjectNamespaceHandler),
			stopper:            make(chan struct{}),
			clientset:          clientset,
			log:                log,
		}, nil
	}
}

func (resourceManagerHandler *ResourceManagerHandler) addNamespaceHandler(objectNamespaceHandler *ObjectNamespaceHandler) {
	resourceManagerHandler.namespaceHandlers[objectNamespaceHandler.namespaceName] = objectNamespaceHandler
}

func (resourceManagerHandler *ResourceManagerHandler) removeNamespaceHandler(namespaceName string) {
	if _, ok := resourceManagerHandler.namespaceHandlers[namespaceName]; ok {
		resourceManagerHandler.namespaceHandlers[namespaceName].Stop()
		delete(resourceManagerHandler.namespaceHandlers, namespaceName)
	}
}

func (resourceManagerHandler *ResourceManagerHandler) Run() error {

	resourceManagerHandler.log.Info(trace(fmt.Sprintf("Namespace selector: <%s> for resourceManagerHandler <%s>", resourceManagerHandler.resourceManager.Spec.NamespaceSelector.String(), resourceManagerHandler.resourceManager.Name)))
	if resourceManagerHandler.resourceManager.Spec.NamespaceSelector != nil {
		resourceManagerHandler.namespacesInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				objectNamespaceHandler, err := NewObjectNamespaceHandler(resourceManagerHandler.resourceManager, obj.(*v1.Namespace).Name, resourceManagerHandler.clientset, resourceManagerHandler.log)
				if err != nil {
					resourceManagerHandler.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s>.", err))
					return
				}
				//if objectNamespaceHandler.terminating {
				//	resourceManagerHandler.log.Info(trace(fmt.Sprintf("Object adding ignored: <%s> Terminating <%b>", objectNamespaceHandler.fullname, objectNamespaceHandler.terminating)))
				//	return
				//}
				resourceManagerHandler.log.Info(trace(fmt.Sprintf("Adding namespace handler: <%s> for resourceManagerHandler <%s>", objectNamespaceHandler.namespaceName, resourceManagerHandler.resourceManager.Name)))
				resourceManagerHandler.addNamespaceHandler(objectNamespaceHandler)
				go objectNamespaceHandler.Run()
				// TODO: handle terminating state  - obj.(*v1.Namespace).Status.Phase == "Terminating"
			},
			//UpdateFunc: func(oldObj interface{}, obj interface{}) {
			//	objectNamespaceHandler, err := NewObjectNamespaceHandler(resourceManagerHandler.resourceManager, obj.(*v1.Namespace).Name, resourceManagerHandler.clientset, resourceManagerHandler.log)
			//	if err != nil {
			//		resourceManagerHandler.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s> for resourceManagerHandler <%s>", err, resourceManagerHandler.resourceManager.Name))
			//		return
			//	}
			//	//if objectNamespaceHandler.terminating {
			//	//	resourceManagerHandler.log.Info(trace(fmt.Sprintf("Object recreating ignored: <%s> Terminating <%b>", objectNamespaceHandler.namespaceName, objectNamespaceHandler.terminating)))
			//	//	return
			//	//}
			//	if reflect.DeepEqual(obj.(*v1.Namespace), oldObj.(*v1.Namespace)) {
			//		resourceManagerHandler.log.Info(trace(fmt.Sprintf("Namespace is not changed <%s> for resourceManagerHandler <%s>. Update Ignored.", obj.(*v1.Namespace).Name, resourceManagerHandler.resourceManager.Name)))
			//		return
			//	}
			//	resourceManagerHandler.log.Info(trace(fmt.Sprintf("Recreating namespace handler: <%s> for resourceManagerHandler <%s>", objectNamespaceHandler.namespaceName, resourceManagerHandler.resourceManager.Name)))
			//	resourceManagerHandler.removeNamespaceHandler(objectNamespaceHandler.namespaceName)
			//	resourceManagerHandler.addNamespaceHandler(objectNamespaceHandler)
			//	go objectNamespaceHandler.Run()
			//},
			DeleteFunc: func(obj interface{}) {
				objectNamespaceHandler, err := NewObjectNamespaceHandler(resourceManagerHandler.resourceManager, obj.(*v1.Namespace).Name, resourceManagerHandler.clientset, resourceManagerHandler.log)
				if err != nil {
					resourceManagerHandler.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s>  for resourceManagerHandler <%s>.", err, resourceManagerHandler.resourceManager.Name))
					return
				}
				resourceManagerHandler.log.Info(trace(fmt.Sprintf("Deleting namespace handler: <%s>", objectNamespaceHandler.namespaceName)))
				resourceManagerHandler.removeNamespaceHandler(objectNamespaceHandler.namespaceName)
			},
		})
		// start the informer
		go resourceManagerHandler.namespacesInformer.Run(resourceManagerHandler.stopper)
	} else {
		objectNamespaceHandler, err := NewObjectNamespaceHandler(resourceManagerHandler.resourceManager, resourceManagerHandler.resourceManager.Namespace, resourceManagerHandler.clientset, resourceManagerHandler.log)
		if err != nil {
			resourceManagerHandler.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s>  for resourceManagerHandler <%s>", err, resourceManagerHandler.resourceManager.Name))
			return err
		}
		//if objectNamespaceHandler.terminating {
		//	resourceManagerHandler.log.Info(trace(fmt.Sprintf("Object adding ignored: <%s> Terminating <%b>", objectNamespaceHandler.fullname, objectNamespaceHandler.terminating)))
		//	return
		//}
		resourceManagerHandler.log.Info(trace(fmt.Sprintf("Adding namespace handler staticly: <%s>  for resourceManagerHandler <%s>", objectNamespaceHandler.namespaceName, resourceManagerHandler.resourceManager.Name)))
		resourceManagerHandler.addNamespaceHandler(objectNamespaceHandler)
		go objectNamespaceHandler.Run()
	}

	return nil
}

func (resourceManagerHandler *ResourceManagerHandler) Stop() {
	// stop channel
	close(resourceManagerHandler.stopper)

	// TODO: cleanup
	for _, h := range resourceManagerHandler.namespaceHandlers {
		h.Stop()
	}
}
