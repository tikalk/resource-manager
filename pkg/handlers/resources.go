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
	clientset   *kubernetes.Clientset
	resource    *v1alpha1.ResourceManager
	objHandlers map[string]*ObjectHandler
}

func NewResourceManagerHandler(r *v1alpha1.ResourceManager) *ResourceManagerHandler {
	// more logic

	return &ResourceManagerHandler{
		resource:    r,
		objHandlers: make(map[string]*ObjectHandler),
	}
}

func (r *ResourceManagerHandler) Start() error {

	// collection := make(map[string]chan struct{})
	// create informer
	selector, _ := metav1.LabelSelectorAsSelector(r.resource.Spec.Selector)

	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})
	factory := informers.NewSharedInformerFactoryWithOptions(r.clientset, 0 /*informers.WithNamespace(p.namespace),*/, labelOptions)
	r.informer = factory.Core().V1().Namespaces().Informer()

	// listen to events
	// create objectHandler
	// ...
}

func (r *ResourceManagerHandler) Stop() error {
	// stop channel
}

type ObjectHandler struct {
}
