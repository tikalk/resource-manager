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

func HandleNamespace(ctx context.Context,
	stopper chan struct{},
	action string,
	condition []resourcemanagmentv1alpha1.ExpiryCondition,
	labelSelector metav1.LabelSelector,
	clientset *kubernetes.Clientset) {
	l := log.FromContext(ctx)
	l.Info(fmt.Sprintf("HandleNamespace begin: action <%s>", action))

	//defer close(stopper)

	select {
	case <-stopper:
		//panic(fmt.Errorf("HandleNamespace: %s", "channel is closed."))
		l.Info(fmt.Sprintf("HandleNamespace abort: action <%s>", action))
		return
	default:
		l.Info(fmt.Sprintf("HandleNamespace: %s", "channel is active"))

	}

	//labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
	//	opts.LabelSelector = labelSelector.String()
	//})
	//factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace("") , labelOptions)
	//
	//labelOptions := informers.WithTweakListOptions(func(opts *metav1.ListOptions) { opts.LabelSelector = labelSelector.String()})

	//factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, labelOptions)
	factory := informers.NewSharedInformerFactory(clientset, time.Second*1)
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

	//informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc:    addFunc,
	//	UpdateFunc: updateFunc,
	//	UpdateFunc: deleteFunc,
	//})

	informer.Run(stopper)

	//// get all the namespaces with the desired selector labels
	//namespaces, err := GetNamespacesByLabel(r, ctx, selectorLabels, l)
	//if err != nil {
	//	l.Error(err, fmt.Sprintf("cannot list namespaces\n"))
	//}
	//
	//if len(namespaces) <= 0 {
	//	fmt.Printf("did not found any namespaces with the requested label\n")
	//	return
	//}
	//
	//fmt.Printf("found %d namespaces with the requested label\n", len(namespaces))
	//
	//for _, namespace := range namespaces {
	//	expired, secondsUntilExpire := utils.IsObjExpired(namespace.CreationTimestamp, condition[0].After)
	//	if expired {
	//		switch action {
	//		case "delete":
	//			l.Info(fmt.Sprintf("namespace '%s' has been expired and will be deleted \n", namespace.Name))
	//			err := r.Delete(ctx, namespace.DeepCopy(), &client.DeleteOptions{})
	//			if err != nil {
	//				l.Error(err, fmt.Sprintf("cannot delete namespaces\n"))
	//			}
	//			l.Info(fmt.Sprintf("namespace '%s' has been deleted \n", namespace.Name))
	//		}
	//	} else {
	//		fmt.Printf("%d seconds has left to namespace '%s' \n", secondsUntilExpire, namespace.Name)
	//
	//	}
	//}

	<-stopper
	l.Info(fmt.Sprintf("HandleNamespace end: action <%s>", action))

}

func addFunc(obj interface{}) {

	//namespaceObj := obj.(*v1.Namespace)
	//
	//// check if the namespace contains "dev"
	//if strings.Contains(namespaceObj.Name, configFile.NamespaceShouldContain) {
	//	fmt.Printf("Found namespace that contains %s: %s \n", configFile.NamespaceShouldContain, namespaceObj.Name)
	//
	//	// find existing Resource Quotas
	//	existingResourceQuotas, _ := GetResourceQuotas(clientset, context.Background(), namespaceObj.Name)
	//
	//	// check if there is any ResourceQuota in namespace
	//	if len(existingResourceQuotas) == 0 {
	//		fmt.Printf("did not found any resource quotas in namespace %s\n", namespaceObj.Name)
	//		CreateCustomResourceQuota(namespaceObj.Name, clientset, configFile)
	//	} else {
	//		// if there is any, check if the mem-cpu-dev-quota existing
	//		for _, quota := range existingResourceQuotas {
	//			if quota == configFile.ResourceQuotaName {
	//				fmt.Printf("ResourceQuota: %s already exists in namespace %s. skipping\n", quota, namespaceObj.Name)
	//			} else {
	//				CreateCustomResourceQuota(namespaceObj.Name, clientset, configFile)
	//			}
	//
	//		}
	//	}
	//
	//}
}

func deleteFunc(obj interface{}) {
	fmt.Println("deleteFunc")
}

func updateFunc(oldObj interface{}, newObj interface{}) {

	switch newObj.(type) {
	//case *v1.Pod:
	//	fmt.Println("updatePod", newObj.(*v1.Pod).Name)
	//case *v1.Service:
	//	fmt.Println("updatePod", newObj.(*v1.Service).Name)
	default:
		fmt.Printf("update: %T\n", newObj)
	}

	//case *v1.Pod:
	//
	//if pod, ok := newObj.(*v1.Pod); ok {
	//	fmt.Println("updatePod", pod.Name)
	//} else if svc, ok := newObj.(*v1.Service); ok {
	//	fmt.Println("updateSvc", svc.Name)
	//} else {
	//	fmt.Println("updateUnknown")
	//}
}
