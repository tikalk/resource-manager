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

	"github.com/tikalk/resource-manager/api/v1alpha1"
	"github.com/tikalk/resource-manager/controllers/handlers"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ResourceManagerReconciler reconciles a ResourceManager object
type ResourceManagerReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	collection map[types.NamespacedName]handlers.Obj
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

func (r *ResourceManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	l := log.FromContext(ctx)
	name := req.NamespacedName.String()

	// your logic here
	resourceManagerObj := &v1alpha1.ResourceManager{}
	err := r.Get(ctx, req.NamespacedName, resourceManagerObj)
	if err != nil {
		if errors.IsNotFound(err) {
			l.Info(fmt.Sprintf("ResourceManager object %s has Not Found!!! \n", req.NamespacedName))
			//r.collection[req.NamespacedName].Stop <- true

			// delete the key from collection map
			delete(r.collection, req.NamespacedName)
			return ctrl.Result{}, nil
		}

		l.Error(err, fmt.Sprintf("Failed reconcile obj %s", name))
	}

	fmt.Printf("found ResourceManager object: %s \n", resourceManagerObj.Name)

	// config handler object
	h := handlers.InitObj(*resourceManagerObj, r.Client, ctx, l)

	//// check if resource exists in our collection, if so, delete
	//if _, ok := r.collection[h.Name]; ok {
	//	l.Info(fmt.Sprintf("Stopping loop for %s\n", h.Name))
	//	r.collection[h.Name]
	//	// delete the key from collection map
	//	delete(r.collection, h.Name)
	//}

	// add handler to collection
	r.collection[req.NamespacedName] = h

	r.collection[req.NamespacedName].Run()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.collection = make(map[types.NamespacedName]handlers.Obj)
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ResourceManager{}).
		Complete(r)
}
