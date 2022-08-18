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
	"time"

	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
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
func (r *ResourceManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	name := req.NamespacedName.String()

	// your logic here
	// Get the resourcemanager Operator object.
	resourceManagerObj := &resourcemanagmentv1alpha1.ResourceManager{}
	err := r.Get(ctx, req.NamespacedName, resourceManagerObj)
	// If one is not found, log a message
	if err != nil && errors.IsNotFound(err) {
		l.Error(err, fmt.Sprintf("Failed reconcile obj %s , Operator resource not found", name))

		// If the Operator is unable to access its custom resource (for any reason besides a simple IsNotFound error)
		// set the OperatorDegraded condition to True with the reason OperatorResourceNotAvailable.
	} else if err != nil {
		l.Error(err, "Error getting operator resource object")
		meta.SetStatusCondition(&resourceManagerObj.Status.Conditions, metav1.Condition{
			Type:               "OperatorDegraded",
			Status:             metav1.ConditionTrue,
			Reason:             "OperatorResourceNotAvailable",
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            fmt.Sprintf("unable to get operator custom resource: %s", err.Error()),
		})
		return ctrl.Result{}, utilerrors.NewAggregate([]error{err, r.Status().Update(ctx, resourceManagerObj)})
	}

	// pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	l.Info(fmt.Sprintf("Starting reconcile obj name: %s, obj: %+v", name, resourceManagerObj))

	deploy := &appsv1.DeploymentList{}
	err = r.Client.List(ctx, deploy, &client.ListOptions{})
	fmt.Printf("There are %d deployments in the cluster\n", len(deploy.Items))

	// Finally, if the Reconcile() function has completed with no critical errors
	// set the OperatorDegraded condition to False with the reason operatorv1alpha1.ReasonSucceeded,
	l.Info(fmt.Sprintf("Done reconcile 12-- obj %s", name))
	meta.SetStatusCondition(&resourceManagerObj.Status.Conditions, metav1.Condition{
		Type:               "OperatorDegraded",
		Status:             metav1.ConditionFalse,
		Reason:             "list deployments succeeded",
		LastTransitionTime: metav1.NewTime(time.Now()),
		Message:            "operator successfully reconciling",
	})
	return ctrl.Result{}, utilerrors.NewAggregate([]error{err, r.Status().Update(ctx, resourceManagerObj)})
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&resourcemanagmentv1alpha1.ResourceManager{}).
		Complete(r)
}
