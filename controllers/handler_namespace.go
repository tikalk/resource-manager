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
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

//+kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=*,resources=namespaces/finalizers,verbs=update

func actionHandler(ctx context.Context, stopper chan struct{}, namespace string,
	clientset *kubernetes.Clientset) {
	l := log.FromContext(ctx)
	l.Info(fmt.Sprintf("actionHandler begin <%s>", namespace))

out:
	for {
		select {
		case <-stopper:
			//panic(fmt.Errorf("HandleNamespace: %s", "channel is closed."))
			//l.Info(fmt.Sprintf("HandleNamespace <%s>: action <%s>", action))
			l.Info(fmt.Sprintf("actionHandler aborted <%s>", namespace))
			break out
		default:

			l.Info(fmt.Sprintf("actionHandler action: <%s> ", namespace))
			// Calling Sleep method
			time.Sleep(5 * time.Second)
		}
	}
	l.Info(fmt.Sprintf("actionHandler end <%s>", namespace))
}

func HandleNamespace(ctx context.Context,
	stopper chan struct{},
	action string,
	condition []resourcemanagmentv1alpha1.ExpiryCondition,
	managedResource resourcemanagmentv1alpha1.ResourceSelector,
	namespace string,
	clientset *kubernetes.Clientset) {

	l := log.FromContext(ctx)
	l.Info(fmt.Sprintf("HandleNamespace begin <%s>", managedResource.Kind))

	//collection := make(map[types.NamespacedName]chan struct{})
	collection := make(map[string]chan struct{})

	selector, _ := metav1.LabelSelectorAsSelector(managedResource.Selector)
	//l.Info(fmt.Sprintf("HandleNamespace selector: <%s>", selector.String()))

	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(namespace), labelOptions)
	informer := factory.Core().V1().Namespaces().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			fmt.Printf("Namespace added: %s\nLabels - %v\n\n", name, labels)
			stopper := make(chan struct{})
			collection[name] = stopper
			go actionHandler(ctx, stopper, name, clientset)
		},

		UpdateFunc: func(oldObj interface{}, obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			fmt.Printf("Namespace updated: %s\nLabels - %v\n\n", name, labels)
			close(collection[name])
			delete(collection, name)
			stopper := make(chan struct{})
			collection[name] = stopper
			go actionHandler(ctx, stopper, name, clientset)
		},

		DeleteFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			fmt.Printf("Namespace deleted: %s\nLabels - %v\n\n", name, labels)
			close(collection[name])
			delete(collection, name)
		},
	})

	informer.Run(stopper)

	l.Info(fmt.Sprintf("HandleNamespace cleanup: action <%s>", action))
	for _, stopper := range collection {
		close(stopper)
	}

	l.Info(fmt.Sprintf("HandleNamespace end: action <%s>", action))

}
