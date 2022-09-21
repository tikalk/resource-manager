package handlers

import (
	v1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceManagerHandler struct {
	informer    cache.SharedIndexInformer
	factory     informers.SharedInformerFactory
	clientset   *kubernetes.Clientset
	resource    *v1alpha1.ResourceManager
	objHandlers map[string]*ObjectHandler
}

func NewResourceManagerHandler(res *v1alpha1.ResourceManager, clientset *kubernetes.Clientset) *ResourceManagerHandler {

	selector, _ := metav1.LabelSelectorAsSelector(res.Spec.Selector)
	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0 /*informers.WithNamespace(p.namespace),*/, labelOptions)

	return &ResourceManagerHandler{
		resource:    res,
		factory:     factory,
		clientset:   clientset,
		objHandlers: make(map[string]*ObjectHandler),
	}
}

func (r *ResourceManagerHandler) getInformer(kind string) (i cache.SharedIndexInformer) {
	switch kind {
	case "Deployment":
		i = r.factory.Apps().V1().Deployments().Informer()
	case "Namespace":
		i = r.factory.Core().V1().Namespaces().Informer()
	}
	return i
}

func (r *ResourceManagerHandler) Start() error {

	// create informer
	r.informer = r.getInformer(r.resource.Spec.ResourceKind)

	// TODO: listen to events
	// r.informer.AddEventHandler()

	// create objectHandler
	// ...
}

func (r *ResourceManagerHandler) Stop() error {
	// stop channel
}

func (r *ResourceManagerHandler) AddFunc() error {

}

type ObjectHandler struct {
}
