package controllers

import (
	"context"
	"fmt"
	"github.com/tikalk/resource-manager/api/v1alpha1"
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

func actionHandler(ctx context.Context, stopper chan struct{}, managedObj *v1.Namespace, managerSpec v1alpha1.ResourceManagerSpec, clientset *kubernetes.Clientset) {
	l := log.FromContext(ctx)
	l.Info(fmt.Sprintf("actionHandler begin <%s>", managedObj.Name))

	var wait time.Duration
	//var duration Duration
	//var age Duration
	switch managerSpec.Condition.Type {
	case "expiry":
		duration, _ := time.ParseDuration(managerSpec.Condition.After)
		age := time.Now().Sub(managedObj.ObjectMeta.CreationTimestamp.Time)
		wait = duration - age

		//waitSec := int64(duration.Seconds()) - (time.Now().Unix() - managedObj.ObjectMeta.CreationTimestamp.Unix())
		l.Info(fmt.Sprintf("actionHandler object expired <%s> duration <%s> age <%s> wait <%s>",
			managedObj.Name,
			duration.String(),
			age.String(),
			wait.String()))
		break
	case "at":
		break
	}

	if wait <= 0 {
		l.Info(fmt.Sprintf("actionHandler object already expired <%s>", managedObj.Name))
	} else {
		select {
		case <-stopper:
			//panic(fmt.Errorf("HandleNamespace: %s", "channel is closed."))
			//l.Info(fmt.Sprintf("HandleNamespace <%s>: action <%s>", action))
			l.Info(fmt.Sprintf("actionHandler aborted <%s>", managedObj.Name))
			return
		case <-time.After(wait):
			l.Info(fmt.Sprintf("actionHandler object expired <%s>", managedObj.Name))
			break
		}
	}

	if managerSpec.DryRun {
		l.Info(fmt.Sprintf("actionHandler dry-run performing object <%s> action <%s> ", managedObj.Name, managerSpec.Action))
	} else {
		l.Info(fmt.Sprintf("actionHandler performing object <%s> action <%s> ", managedObj.Name, managerSpec.Action))
	}
	//	// Calling Sleep method
	//	time.Sleep(5 * time.Second)

	//}

	l.Info(fmt.Sprintf("actionHandler end <%s>", managedObj.Name))
}

func HandleNamespace(p HandlerParams) {

	//	ctx context.Context,
	//	stopper chan struct{},
	//	action string,
	//	condition []resourcemanagmentv1alpha1.ExpiryCondition,
	//	resourceKind string,
	//	labelSelector *metav1.LabelSelector,
	//	namespace string,
	//	clientset *kubernetes.Clientset
	//
	//)

	l := log.FromContext(p.ctx)
	l.Info(fmt.Sprintf("HandleNamespace begin <%s>", p.name))

	//collection := make(map[types.NamespacedName]chan struct{})
	collection := make(map[string]chan struct{})

	selector, _ := metav1.LabelSelectorAsSelector(p.spec.Selector)

	labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = selector.String()
	})

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(p.namespace), labelOptions)
	informer := factory.Core().V1().Namespaces().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			l.Info(fmt.Sprintf("Namespace added: <%s> Created at <%s>", name, obj.(*v1.Namespace).ObjectMeta.CreationTimestamp.String()))
			stopper := make(chan struct{})
			collection[name] = stopper
			go actionHandler(p.ctx, stopper, obj.(*v1.Namespace), p.spec, clientset)
		},

		UpdateFunc: func(oldObj interface{}, obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			l.Info(fmt.Sprintf("Namespace updated: %s\nLabels - %v\n\n", name, labels))
			close(collection[name])
			delete(collection, name)
			stopper := make(chan struct{})
			collection[name] = stopper
			go actionHandler(p.ctx, stopper, obj.(*v1.Namespace), p.spec, clientset)
		},

		DeleteFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			labels := obj.(*v1.Namespace).Labels
			l.Info(fmt.Sprintf("Namespace deleted: %s\nLabels - %v\n\n", name, labels))
			close(collection[name])
			delete(collection, name)
		},
	})

	informer.Run(p.stopper)

	l.Info(fmt.Sprintf("HandleNamespace cleanup <%s>", p.name))
	for _, stopper := range collection {
		close(stopper)
	}

	l.Info(fmt.Sprintf("HandleNamespace end <%s>", p.name))

}
