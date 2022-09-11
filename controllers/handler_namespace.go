package controllers

import (
	"fmt"
	"github.com/tikalk/resource-manager/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"time"
)

//+kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=*,resources=namespaces/finalizers,verbs=update

func namespaceActionHandler(stopper chan struct{}, managedObj *v1.Namespace, mgrSpec v1alpha1.ResourceManagerSpec) {

	l.Info(fmt.Sprintf("namespaceActionHandler begin <%s>", managedObj.Name))

	var wait time.Duration
	//var duration Duration
	//var age Duration
	switch mgrSpec.Condition.Type {
	case "expiry":
		duration, _ := time.ParseDuration(mgrSpec.Condition.After)
		age := time.Now().Sub(managedObj.ObjectMeta.CreationTimestamp.Time)
		wait = duration - age

		l.Info(trace(fmt.Sprintf("object expired <%s> duration <%s> age <%s> wait <%s>",
			managedObj.Name,
			duration.String(),
			age.String(),
			wait.String())))
		break
	case "at":
		break
	}

	if wait <= 0 {
		l.Info(trace(fmt.Sprintf("object already expired <%s>", managedObj.Name)))
	} else {
		select {
		case <-stopper:
			l.Info(trace(fmt.Sprintf("aborted <%s>", managedObj.Name)))
			return
		case <-time.After(wait):
			l.Info(trace(fmt.Sprintf("object expired <%s>", managedObj.Name)))
			break
		}
	}

	if mgrSpec.DryRun {
		l.Info(trace(fmt.Sprintf("dry-run performing object <%s> action <%s> ", managedObj.Name, mgrSpec.Action)))
	} else {
		l.Info(trace(fmt.Sprintf("performing object <%s> action <%s> ", managedObj.Name, mgrSpec.Action)))
	}

	l.Info(trace(fmt.Sprintf("end <%s>", managedObj.Name)))
}

func HandleNamespace(p HandlerParams) {

	l.Info(trace(fmt.Sprintf("begin <%s>", p.name)))

	//collection := make(map[types.NamespacedName]chan struct{})
	collection := make(map[string]chan struct{})

	selector, _ := metav1.LabelSelectorAsSelector(p.mgrSpec.Selector)

	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0 /*informers.WithNamespace(p.namespace),*/, labelOptions)
	informer := factory.Core().V1().Namespaces().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			l.Info(trace(fmt.Sprintf("Namespace added: <%s> Created at <%s>", name, obj.(*v1.Namespace).ObjectMeta.CreationTimestamp.String())))
			stopper := make(chan struct{})
			collection[name] = stopper
			go namespaceActionHandler(stopper, obj.(*v1.Namespace), p.mgrSpec)
		},

		UpdateFunc: func(oldObj interface{}, obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			l.Info(trace(fmt.Sprintf("Namespace updated: %s", name)))
			close(collection[name])
			delete(collection, name)
			stopper := make(chan struct{})
			collection[name] = stopper
			go namespaceActionHandler(stopper, obj.(*v1.Namespace), p.mgrSpec)
		},

		DeleteFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			l.Info(trace(fmt.Sprintf("Namespace deleted: %s\nLabels - %v\n\n", name, labels)))
			close(collection[name])
			delete(collection, name)
		},
	})

	informer.Run(p.stopper)

	l.Info(trace(fmt.Sprintf("cleanuping... <%s>", p.name)))
	for _, stopper := range collection {
		close(stopper)
	}

	l.Info(trace(fmt.Sprintf("end <%s>", p.name)))

}
