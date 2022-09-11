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
	"github.com/go-logr/logr"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ResourceManagerReconciler reconciles a ResourceManager object
type ResourceManagerReconciler struct {
	client.Client
	Scheme *k8sruntime.Scheme
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

var clientset *kubernetes.Clientset

type HandlerParams struct {
	stopper   chan struct{}
	mgrSpec   resourcemanagmentv1alpha1.ResourceManagerSpec
	name      string
	namespace string
}

func handlerFactory(resourceKind string) (func(HandlerParams), error) {

	switch resourceKind {
	case "Namespace":
		return HandleNamespace, nil
	default:
		return nil, fmt.Errorf("unexpected resourceKind <%s>", resourceKind)
	}

}

var l logr.Logger
var loggerInitialized bool

func (r *ResourceManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	if !loggerInitialized {
		l = log.FromContext(ctx)
		loggerInitialized = true
	}

	//name := req.NamespacedName.String()
	l.Info(trace(fmt.Sprintf("ResourceManager object %s reconciled. Reconciling...", req.NamespacedName)))

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
		return ctrl.Result{}, nil
	}

	if _, ok := collection[req.NamespacedName]; ok {
		l.Info(fmt.Sprintf("ResourceManager %s changed. Recreating...", req.NamespacedName))
		close(collection[req.NamespacedName])
		delete(collection, req.NamespacedName)
	} else {
		l.Info(fmt.Sprintf("ResourceManager %s created. Creating...", req.NamespacedName))
	}

	collection[req.NamespacedName] = make(chan struct{})

	if !resourceManager.Spec.Disabled {
		handler, err := handlerFactory(resourceManager.Spec.ResourceKind)
		if err != nil {
			l.Error(err, fmt.Sprintf("Failed to creatate handler for <%s>", req.NamespacedName))
		} else {
			l.Info(fmt.Sprintf("Starting handler %s for ...", req.NamespacedName))
			go handler(HandlerParams{
				collection[req.NamespacedName],
				resourceManager.Spec,
				req.NamespacedName.String(),
				req.NamespacedName.Namespace})
		}
	}

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

func trace(msg string) string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Sprintf("%s:%d %s | %s", "?", 0, "?", msg)
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s:%d %s | %s", file, line, "?", msg)
	}

	return fmt.Sprintf("%s:%d %s | %s", file, line, fn.Name(), msg)
}
