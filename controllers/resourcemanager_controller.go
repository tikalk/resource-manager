/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ResourceManagerReconciler reconciles a ResourceManager object

type ResourceManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=resource-management.tikalk.com,resources=resourcemanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=resource-management.tikalk.com,resources=resourcemanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=resource-management.tikalk.com,resources=resourcemanagers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ResourceManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile

//func HandleNamespace(ctx context.Context,
//	action string,
//	condition []resourcemanagmentv1alpha1.ExpiryCondition,
//	selectorLabels metav1.LabelSelector) {
//	l := log.FromContext(ctx)
//	l.Info(fmt.Sprintf("HandleNamespace: %s", action))
//	//// get all the namespaces with the desired selector labels
//	//namespaces, err := GetNamespacesByLabel(r, ctx, selectorLabels, l)
//	//if err != nil {
//	//	l.Error(err, fmt.Sprintf("cannot list namespaces\n"))
//	//}
//	//
//	//if len(namespaces) <= 0 {
//	//	fmt.Printf("did not found any namespaces with the requested label\n")
//	//	return
//	//}
//	//
//	//fmt.Printf("found %d namespaces with the requested label\n", len(namespaces))
//	//
//	//for _, namespace := range namespaces {
//	//	expired, secondsUntilExpire := utils.IsObjExpired(namespace.CreationTimestamp, condition[0].After)
//	//	if expired {
//	//		switch action {
//	//		case "delete":
//	//			l.Info(fmt.Sprintf("namespace '%s' has been expired and will be deleted \n", namespace.Name))
//	//			err := r.Delete(ctx, namespace.DeepCopy(), &client.DeleteOptions{})
//	//			if err != nil {
//	//				l.Error(err, fmt.Sprintf("cannot delete namespaces\n"))
//	//			}
//	//			l.Info(fmt.Sprintf("namespace '%s' has been deleted \n", namespace.Name))
//	//		}
//	//	} else {
//	//		fmt.Printf("%d seconds has left to namespace '%s' \n", secondsUntilExpire, namespace.Name)
//	//
//	//	}
//	//}
//}

var collection map[types.NamespacedName]chan struct{}

// var sharedInformerFactory informers.SharedInformerFactory
var clientset *kubernetes.Clientset

//type handlerFunc func(resourcemanagmentv1alpha1.ResourceManager)

func handlerFactory(ctx context.Context, resourceType string) func(
	ctx context.Context,
	stopper chan struct{},
	action string,
	condition []resourcemanagmentv1alpha1.ExpiryCondition,
	selectorLabels metav1.LabelSelector,
	clientset *kubernetes.Clientset) {

	return func(ctx context.Context,
		stopper chan struct{},
		action string,
		condition []resourcemanagmentv1alpha1.ExpiryCondition,
		selector metav1.LabelSelector,
		clientset *kubernetes.Clientset) {
		for {
			switch resourceType {
			case "Namespace":
				HandleNamespace(ctx, stopper, action, condition, selector, clientset)
			}
			//select {
			//case <-stopper:
			//	return
			//default:
			//	time.Sleep(time.Second * 5)
			//}
		}
	}
}

func (r *ResourceManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	name := req.NamespacedName.String()

	l := log.FromContext(ctx)

	l.Info(fmt.Sprint("Reconcile begin ", req))

	l.Info(fmt.Sprintf("Collection begin"))
	for i := range collection {
		l.Info(fmt.Sprintf("Collection key: %s", i.String()))
	}
	l.Info(fmt.Sprintf("Collection end"))

	l.Info(fmt.Sprintf("Resource manager: %s", name))

	// your logic here
	stopper, chExist := collection[req.NamespacedName]

	resourceManager := resourcemanagmentv1alpha1.ResourceManager{}
	if r.Get(ctx, req.NamespacedName, &resourceManager) != nil {
		l.Info(fmt.Sprintf("Resource manager deleted: %s", name))
		if chExist {
			close(stopper)
			delete(collection, req.NamespacedName)
		}
		return ctrl.Result{}, nil
	}

	// Remove object from collection
	if chExist {
		close(stopper)
		delete(collection, req.NamespacedName)
	} else {

	}

	// Check enabling flag and skip if disabled
	if !resourceManager.Spec.Active {
		return ctrl.Result{}, nil
	}

	// Create new object by object factory
	if !chExist {
		stopper = make(chan struct{})
	}
	collection[req.NamespacedName] = stopper
	handler := handlerFactory(ctx, resourceManager.Spec.Resources)
	go handler(ctx, stopper, resourceManager.Spec.Action, resourceManager.Spec.Condition, *resourceManager.Spec.Selector, clientset)

	//l.Info(fmt.Sprintf("Done reconcile 12-- obj %s", collection))

	// pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	//l.Info(fmt.Sprintf("Starting reconcile obj name: %s, obj: %+v", name, resourceManagerObj))

	//collection[name] :=
	//chan bool
	//go func(stop chan bool) {
	//	for {
	//		handlers.HandleNamespaceObj(r.Client, ctx, selectorLabels, action, condition, l)
	//		select {
	//		case <-stop:
	//			return
	//		default:
	//			time.Sleep(time.Second)
	//		}
	//	}
	//}

	//deploy := &appsv1.DeploymentList{}
	//err = r.Client.List(ctx, deploy, &client.ListOptions{})
	//fmt.Printf("There are %d deployments in the cluster\n", len(deploy.Items))

	//namespaceHandler := func(stopper chan struct{}) {
	//	for {
	//		handlers.HandleNamespaceObj(r.Client, ctx, selectorLabels, action, condition, l)
	//		select {
	//		case <-stop:
	//			return
	//		default:
	//			time.Sleep(time.Second)
	//		}
	//	}
	//}

	//handlerFactory()
	//switch ResourceType {
	//case "namespace":
	//	if _, ok := collection[req.NamespacedName]; ok {
	//		collection[req.NamespacedName].stop <- true
	//	}
	//	collection[req.NamespacedName] = &Handler{
	//		F: namenamespaceHandler,
	//		Stop: chan bool
	//	}
	//	h = collection[req.NamespacedName]
	//	go h.F(h.Stop)
	//}

	l.Info(fmt.Sprintf("Reconcile end <%s>", name))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	collection = make(map[types.NamespacedName]chan struct{})
	//sharedInformerFactory = informers.NewSharedInformerFactory(r.Client, time.Second*1)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&resourcemanagmentv1alpha1.ResourceManager{}).
		Complete(r)
}
