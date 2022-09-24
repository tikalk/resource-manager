package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/tikalk/resource-manager/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"time"
	//"errors"
)

//+kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=*,resources=namespaces/finalizers,verbs=update

func namespaceActionHandler(stopper chan struct{}, managedObj *v1.Namespace, mgrSpec v1alpha1.ResourceManagerSpec) {

	//l.Info(fmt.Sprintf("namespaceActionHandler begin <%s>", managedObj.Name))

	var wait time.Duration
	//var duration Duration
	//var age Duration
	if mgrSpec.Condition.Timeframe != "" {
		timeframe, _ := time.ParseDuration(mgrSpec.Condition.Timeframe)
		age := time.Now().Sub(managedObj.ObjectMeta.CreationTimestamp.Time)
		wait = timeframe - age

		l.Info(trace(fmt.Sprintf("object timeframe expiration <%s> timeframe <%s> age <%s> wait <%s>",
			managedObj.Name,
			timeframe.String(),
			age.String(),
			wait.String())))
	} else if mgrSpec.Condition.ExpireAt != "" {
		expireAt, err := time.Parse("15:04", mgrSpec.Condition.ExpireAt)
		if err != nil {
			l.Error(err, fmt.Sprintf("Failed to parse %s", mgrSpec.Condition.ExpireAt))
		}

		now := time.Now()

		if expireAt.Hour()*60+expireAt.Minute() > now.Hour()*60+now.Minute() {
			wait = time.Date(now.Year(), now.Month(), now.Day(), expireAt.Hour(), expireAt.Minute(), 0, 0, now.Location()).Sub(now)
		} else {
			tomorrow := now.Add(24 * time.Hour)
			wait = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), expireAt.Hour(), expireAt.Minute(), 0, 0, tomorrow.Location()).Sub(now)
		}

		l.Info(trace(fmt.Sprintf("object time expiration <%s> expireAt <%s> now <%s> wait <%s>",
			managedObj.Name,
			expireAt.String(),
			now.String(),
			wait)))
	} else {
		l.Error(errors.New("expiration is not configured"), trace(fmt.Sprintf("object hanler <%s> aborted", managedObj.Name)))
		return
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
		var err error
		switch mgrSpec.Action {
		case "delete":
			err = namespaceDelete(managedObj.Name)
			break
		default:
			err = errors.New("unexpected action")
		}
		if err != nil {
			l.Error(err, trace(fmt.Sprintf("object <%s> action <%s> failed", managedObj.Name, mgrSpec.Action)))
		} else {
			l.Info(trace(fmt.Sprintf("object <%s> action <%s> finished", managedObj.Name, mgrSpec.Action)))
		}

	}

	//l.Info(trace(fmt.Sprintf("end <%s>", managedObj.Name)))
}

func namespaceDelete(objName string) error {
	var opts metav1.DeleteOptions
	ctx := context.TODO()
	return clientset.CoreV1().Namespaces().Delete(ctx, objName, opts)
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
			stopper := make(chan struct{})
			collection[name] = stopper
			if obj.(*v1.Namespace).Status.Phase == "Terminating" {
				l.Info(trace(fmt.Sprintf("Namespace add ignored: <%s> Status <%s>", name, obj.(*v1.Namespace).Status.Phase)))
			} else {
				l.Info(trace(fmt.Sprintf("Namespace added: <%s> Created at <%s>", name, obj.(*v1.Namespace).ObjectMeta.CreationTimestamp.String())))
				go namespaceActionHandler(stopper, obj.(*v1.Namespace), p.mgrSpec)
			}

		},

		UpdateFunc: func(oldObj interface{}, obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			close(collection[name])
			delete(collection, name)
			stopper := make(chan struct{})
			collection[name] = stopper
			if obj.(*v1.Namespace).Status.Phase == "Terminating" {
				l.Info(trace(fmt.Sprintf("Namespace update ignored: <%s> Status <%s>", name, obj.(*v1.Namespace).Status.Phase)))
			} else {
				l.Info(trace(fmt.Sprintf("Namespace updated: %s", name)))
				go namespaceActionHandler(stopper, obj.(*v1.Namespace), p.mgrSpec)
			}
		},

		DeleteFunc: func(obj interface{}) {
			//fmt.Println("namespace add")
			name := obj.(*v1.Namespace).Name
			//labels := obj.(*v1.Namespace).Labels
			l.Info(trace(fmt.Sprintf("Namespace deleted: <%s>", name)))
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
