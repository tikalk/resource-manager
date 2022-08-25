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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
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

//+kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ResourceManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile

var collection map[types.NamespacedName]chan struct{}

//type FHandler struct {
//	F    func(ctx context.Context,
//		//stopper chan struct{},
//		stop chan bool,
//		action string,
//		condition []resourcemanagmentv1alpha1.ExpiryCondition,
//		selector metav1.LabelSelector,
//		clientset *kubernetes.Clientset)
//	Stop chan bool
//}

var clientset *kubernetes.Clientset

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
		l := log.FromContext(ctx)
		switch resourceType {
		case "namespace":
			HandleNamespace(ctx, stopper, action, condition, selector, clientset)
		default:
			l.Error(errors.NewBadRequest("math: square root of negative number"), fmt.Sprintf("Unexpected resourceType %s", resourceType))
		}
		//select {
		//case <-stopper:
		//	return
		//default:
		//	time.Sleep(time.Second * 5)
		//}
	}
}

func (r *ResourceManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	l := log.FromContext(ctx)
	//name := req.NamespacedName.String()
	l.Info(fmt.Sprintf("ResourceManager object %s reconciled. Reconciling...", req.NamespacedName))

	// your logic here
	resourceManager := &resourcemanagmentv1alpha1.ResourceManager{}
	err := r.Get(ctx, req.NamespacedName, resourceManager)
	if err != nil {
		if errors.IsNotFound(err) {
			l.Info(fmt.Sprintf("ResourceManager object %s deleted. Removing...", req.NamespacedName))
			close(collection[req.NamespacedName])
			delete(collection, req.NamespacedName)
			return ctrl.Result{}, nil
		}

		l.Error(err, fmt.Sprintf("Failed reconcile obj %s", req.NamespacedName))
	}

	fmt.Printf("found ResourceManager object: %s \n", resourceManager.Name)

	//// config handler object
	//h := handlers.Obj{
	//	Name: req.NamespacedName,
	//	C:    r.Client,
	//	Ctx:  ctx,
	//	L:    l,
	//	Spec: resourceManagerObj.Spec,
	//}
	//
	//l.Info(fmt.Sprintf(
	//	"\n"+
	//		" ResourceType: %s \n"+
	//		" selectorLables %s \n"+
	//		" action: %s \n"+
	//		" condition: %s \n"+
	//		" type: %s \n",
	//	h.Spec.Resources,
	//	h.Spec.Selector.MatchLabels,
	//	h.Spec.Action,
	//	h.Spec.Condition[0].After,
	//	h.Spec.Condition[0].Type))

	// check if resource exists in our collection, if so, delete
	if _, ok := collection[req.NamespacedName]; ok {
		l.Info(fmt.Sprintf("ResourceManager object %s found. Do nothing...", req.NamespacedName))
		return ctrl.Result{}, nil
	}

	collection[req.NamespacedName] = make(chan struct{})
	handler := handlerFactory(ctx, resourceManager.Spec.Resources)
	l.Info(fmt.Sprintf("ResourceManager object %s created. Starting handler...", req.NamespacedName))
	go handler(ctx, collection[req.NamespacedName], resourceManager.Spec.Action, resourceManager.Spec.Condition, *resourceManager.Spec.Selector, clientset)

	//
	//deploy := &appsv1.DeploymentList{}
	//err = r.Client.List(ctx, deploy, &client.ListOptions{})
	//fmt.Printf("There are %d deployments in the cluster\n", len(deploy.Items))
	//
	//l.Info(fmt.Sprintf("Done reconcile 12-- obj %s", name))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	collection = make(map[types.NamespacedName]chan struct{})

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&resourcemanagmentv1alpha1.ResourceManager{}).
		Complete(r)
}
