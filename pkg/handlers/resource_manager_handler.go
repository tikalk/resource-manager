package handlers

import (
	"fmt"

	v1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type ResourceManagerHandler struct {
	informer        cache.SharedIndexInformer
	factory         informers.SharedInformerFactory
	clientset       *kubernetes.Clientset
	resourceManager *v1alpha1.ResourceManager
	objHandlers     map[types.NamespacedName]*ObjectHandler

	stopper chan struct{}
}

func NewResourceManagerHandler(res *v1alpha1.ResourceManager, clientset *kubernetes.Clientset) (*ResourceManagerHandler, error) {
	selector, _ := metav1.LabelSelectorAsSelector(res.Spec.Selector)
	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0 /*informers.WithNamespace(p.namespace),*/, labelOptions)

	i, err := createInformer(factory, res.Spec.ResourceKind)
	if err != nil {
		return nil, err
	}

	return &ResourceManagerHandler{
		resourceManager: res,
		factory:         factory,
		informer:        i,
		clientset:       clientset,
		objHandlers:     make(map[types.NamespacedName]*ObjectHandler),
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

func (r *ResourceManagerHandler) addObjHandelr(fullname types.NamespacedName, objHandler *ObjectHandler) {
	r.objHandlers[fullname] = objHandler
}

func (r *ResourceManagerHandler) removeObjHandelr(fullname types.NamespacedName) {
	if _, ok := r.objHandlers[fullname]; ok {
		r.objHandlers[fullname].Stop()
		delete(r.objHandlers, fullname)
	}
}

func (r *ResourceManagerHandler) Start() error {

	// TODO: listen to events
	r.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			objHandler, fullname, _ := NewObjectHandler(r.resourceManager, obj)
			r.addObjHandelr(fullname, objHandler)
			go objHandler.Run()
			// TODO: handle terminating state  - obj.(*v1.Namespace).Status.Phase == "Terminating"
		},
		UpdateFunc: func(oldObj interface{}, obj interface{}) {
			objHandler, fullname, _ := NewObjectHandler(r.resourceManager, obj)
			r.removeObjHandelr(fullname)
			r.addObjHandelr(fullname, objHandler)
			go objHandler.Run()
		},
		DeleteFunc: func(obj interface{}) {
			_, fullname, _ := NewObjectHandler(r.resourceManager, obj)
			r.removeObjHandelr(fullname)
		},
	})
	// start the informer
	go r.informer.Run(r.stopper)

	return nil
}

func (r *ResourceManagerHandler) Stop() {
	// stop channel
	close(r.stopper)

	// TODO: cleanup
	for _, h := range r.objHandlers {
		h.Stop()
	}
}
