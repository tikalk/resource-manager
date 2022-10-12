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

func NewResourceManagerHandler(rm *v1alpha1.ResourceManager, clientset *kubernetes.Clientset, log logr.Logger) (*ResourceManagerHandler, error) {
	if rm.Spec.NamespaceSelector != nil {
		selector, _ := metav1.LabelSelectorAsSelector(rm.Spec.NamespaceSelector)
		labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
			opts.LabelSelector = selector.String()
		})
		factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, labelOptions)

		namespacesInformer := factory.Core().V1().Namespaces().Informer()

		return &ResourceManagerHandler{
			resourceManager:    rm,
			namespacesInformer: namespacesInformer,
			namespaceHandlers:  make(map[string]*ObjectNamespaceHandler),
			stopper:            make(chan struct{}),
			clientset:          clientset,
			log:                log,
		}, nil
	} else {
		return &ResourceManagerHandler{
			resourceManager:    rm,
			namespacesInformer: nil,
			namespaceHandlers:  make(map[string]*ObjectNamespaceHandler),
			stopper:            make(chan struct{}),
			clientset:          clientset,
			log:                log,
		}, nil
	}
}

func (h *ResourceManagerHandler) addNamespaceHandler(objectNamespaceHandler *ObjectNamespaceHandler) {
	h.namespaceHandlers[objectNamespaceHandler.namespaceName] = objectNamespaceHandler
}

func (h *ResourceManagerHandler) removeNamespaceHandler(namespaceName string) {
	if _, ok := h.namespaceHandlers[namespaceName]; ok {
		h.namespaceHandlers[namespaceName].Stop()
		delete(h.namespaceHandlers, namespaceName)
	}
}

func (h *ResourceManagerHandler) Run() error {

	h.log.Info(trace(fmt.Sprintf("Namespace selector: <%s> for h <%s>", h.resourceManager.Spec.NamespaceSelector.String(), h.resourceManager.Name)))
	if h.resourceManager.Spec.NamespaceSelector != nil {
		h.namespacesInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				objectNamespaceHandler, err := NewObjectNamespaceHandler(h.resourceManager, obj.(*v1.Namespace).Name, h.clientset, h.log)
				if err != nil {
					h.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s>.", err))
					return
				}
				//if objectNamespaceHandler.terminating {
				//	h.log.Info(trace(fmt.Sprintf("Object adding ignored: <%s> Terminating <%b>", objectNamespaceHandler.fullname, objectNamespaceHandler.terminating)))
				//	return
				//}
				h.log.Info(trace(fmt.Sprintf("Adding namespace handler: <%s> for h <%s>", objectNamespaceHandler.namespaceName, h.resourceManager.Name)))
				h.addNamespaceHandler(objectNamespaceHandler)
				go objectNamespaceHandler.Run()
				// TODO: handle terminating state  - obj.(*v1.Namespace).Status.Phase == "Terminating"
			},
			//UpdateFunc: func(oldObj interface{}, obj interface{}) {
			//	objectNamespaceHandler, err := NewObjectNamespaceHandler(h.resourceManager, obj.(*v1.Namespace).Name, h.clientset, h.log)
			//	if err != nil {
			//		h.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s> for h <%s>", err, h.resourceManager.Name))
			//		return
			//	}
			//	//if objectNamespaceHandler.terminating {
			//	//	h.log.Info(trace(fmt.Sprintf("Object recreating ignored: <%s> Terminating <%b>", objectNamespaceHandler.namespaceName, objectNamespaceHandler.terminating)))
			//	//	return
			//	//}
			//	if reflect.DeepEqual(obj.(*v1.Namespace), oldObj.(*v1.Namespace)) {
			//		h.log.Info(trace(fmt.Sprintf("Namespace is not changed <%s> for h <%s>. Update Ignored.", obj.(*v1.Namespace).Name, h.resourceManager.Name)))
			//		return
			//	}
			//	h.log.Info(trace(fmt.Sprintf("Recreating namespace handler: <%s> for h <%s>", objectNamespaceHandler.namespaceName, h.resourceManager.Name)))
			//	h.removeNamespaceHandler(objectNamespaceHandler.namespaceName)
			//	h.addNamespaceHandler(objectNamespaceHandler)
			//	go objectNamespaceHandler.Run()
			//},
			DeleteFunc: func(obj interface{}) {
				onh, err := NewObjectNamespaceHandler(h.resourceManager, obj.(*v1.Namespace).Name, h.clientset, h.log)
				if err != nil {
					h.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s>  for h <%s>.", err, h.resourceManager.Name))
					return
				}
				h.log.Info(trace(fmt.Sprintf("Deleting namespace handler: <%s>", onh.namespaceName)))
				h.removeNamespaceHandler(onh.namespaceName)
			},
		})
		// start the informer
		go h.namespacesInformer.Run(h.stopper)
	} else {
		onh, err := NewObjectNamespaceHandler(h.resourceManager, h.resourceManager.Namespace, h.clientset, h.log)
		if err != nil {
			h.log.Error(err, fmt.Sprintf("NewObjectNamespaceHandler handler creating failed with error <%s>  for h <%s>", err, h.resourceManager.Name))
			return err
		}
		//if onh.terminating {
		//	h.log.Info(trace(fmt.Sprintf("Object adding ignored: <%s> Terminating <%b>", onh.fullname, onh.terminating)))
		//	return
		//}
		h.log.Info(trace(fmt.Sprintf("Adding namespace handler staticly: <%s>  for h <%s>", onh.namespaceName, h.resourceManager.Name)))
		h.addNamespaceHandler(onh)
		go onh.Run()
	}

	return nil
}

func (h *ResourceManagerHandler) Stop() {
	// stop channel
	close(h.stopper)

	// TODO: cleanup
	for _, h := range h.namespaceHandlers {
		h.Stop()
	}
}
