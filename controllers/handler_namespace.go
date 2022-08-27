package controllers

import (
	"context"
	"fmt"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

//+kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=*,resources=namespaces/finalizers,verbs=update

func HandleNamespace(ctx context.Context,
	stopper chan struct{},
	action string,
	condition []resourcemanagmentv1alpha1.ExpiryCondition,
	managedResource resourcemanagmentv1alpha1.ResourceSelector,
	namespace string,
	clientset *kubernetes.Clientset) {

	//l := log.FromContext(ctx)
	//l.Info(fmt.Sprintf("HandleNamespace <%s> begin", managedResource.Kind))

	//defer close(stopper)

	//select {
	//case <-stopper:
	//	//panic(fmt.Errorf("HandleNamespace: %s", "channel is closed."))
	//	//l.Info(fmt.Sprintf("HandleNamespace <%s>: action <%s>", action))
	//	return
	//default:
	//	//l.Info(fmt.Sprintf("HandleNamespace: %s", "channel is active"))
	//
	//}

	//l.Info(fmt.Sprintf("HandleNamespace labelSelector: <%s>", managedResource.Selector.String()))

	selector, _ := metav1.LabelSelectorAsSelector(managedResource.Selector)
	//l.Info(fmt.Sprintf("HandleNamespace selector: <%s>", selector.String()))

	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(namespace), labelOptions)

	//labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) { opts.LabelSelector = labelSelector.String()})

	//factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, labelOptions)
	//factory := informers.NewSharedInformerFactory(clientset, time.Second*1)
	informer := factory.Core().V1().Namespaces().Informer()
	//stopper := make(chan struct{})

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			fmt.Printf("Namespace added: %s\nLabels - %v\n\n", name, labels)
		},

		//UpdateFunc: func(oldObj interface{}, obj interface{}) {
		//	//fmt.Println("namespace add")
		//	name := obj.(*v1.Namespace).Name
		//	labels := obj.(*v1.Namespace).Labels
		//	fmt.Printf("Namespace updated: %s\nLabels - %v\n\n", name, labels)
		//},

		DeleteFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			fmt.Printf("Namespace deleted: %s\nLabels - %v\n\n", name, labels)
		},
	})

	informer.Run(stopper)

	//<-stopper
	//l.Info(fmt.Sprintf("HandleNamespace end: action <%s>", action))

}
