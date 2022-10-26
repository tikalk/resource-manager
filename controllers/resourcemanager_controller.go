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
	"k8s.io/apimachinery/pkg/api/errors"
	"os"
	"runtime"
	"time"

	"github.com/go-logr/logr"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"go.uber.org/zap/zapcore"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	uzap "go.uber.org/zap"
	//"go.uber.org/zap/zapcore"
	//logf "sigs.k8s.io/controller-runtime/pkg/log"
	//"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"reflect"
)

//+kubebuilder:rbac:groups=resource-management.tikalk.com,resources=resourcemanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=resource-management.tikalk.com,resources=resourcemanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=resource-management.tikalk.com,resources=resourcemanagers/finalizers,verbs=update

//+kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// ResourceManagerReconciler reconciles a ResourceManager object
type ResourceManagerReconciler struct {
	client.Client
	Scheme                  *k8sruntime.Scheme
	resourceManagerHandlers map[types.NamespacedName]*ResourceManagerHandler

	clientset *kubernetes.Clientset
	log       logr.Logger
}

// registerAndRunResourceManagerHandler add the handler to the collection and then run it
func (r *ResourceManagerReconciler) registerAndRunResourceManagerHandler(resourceManagerName types.NamespacedName, resourceManagerHandler *ResourceManagerHandler) {
	r.resourceManagerHandlers[resourceManagerName] = resourceManagerHandler
	go resourceManagerHandler.Run()

}

// findResourceManagerHandler will look for the resource manager handler object in the collection
func (r *ResourceManagerReconciler) findResourceManagerHandler(resourceManagerName types.NamespacedName) *ResourceManagerHandler {
	return r.resourceManagerHandlers[resourceManagerName]
}

func (r *ResourceManagerReconciler) removeResourceManagerHandler(resourceManagerName types.NamespacedName) {
	if _, ok := r.resourceManagerHandlers[resourceManagerName]; ok {
		// l.Info(fmt.Sprintf("ResourceManager %s changed. Recreating...", req.NamespacedName))
		r.resourceManagerHandlers[resourceManagerName].Stop()
		delete(r.resourceManagerHandlers, resourceManagerName)
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ResourceManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile

// var ctx context.Context
// var initialized bool
func (r *ResourceManagerReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {

	r.log.Info(trace(fmt.Sprintf("ResourceManager object <%s> reconciled. Reconciling...", request.NamespacedName)))

	resourceManager := &resourcemanagmentv1alpha1.ResourceManager{}
	if err := r.Get(ctx, request.NamespacedName, resourceManager); err != nil {
		if errors.IsNotFound(err) {
			r.log.Info(fmt.Sprintf("ResourceManager object %s deleted. Removing...", request.NamespacedName))
			r.removeResourceManagerHandler(request.NamespacedName)
			return ctrl.Result{}, nil
		}

		r.log.Error(err, fmt.Sprintf("Failed reconcile ResourceManager object %s", request.NamespacedName))
		return ctrl.Result{}, nil
	}

	if resourceManagerHandler := r.findResourceManagerHandler(request.NamespacedName); resourceManagerHandler != nil {
		//r.log.Info(trace(fmt.Sprintf("ResourceManager object updated: \nold <%+v> \nnew <%+v>.", oldObj.resourceManager, resourceManager)))
		if reflect.DeepEqual(resourceManager.Spec, resourceManagerHandler.resourceManager.Spec) {
			r.log.Info(trace(fmt.Sprintf("ResourceManager spec is not changed <%s>. Ignoring...", request.NamespacedName)))
			return ctrl.Result{}, nil
		}
		r.log.Info(trace(fmt.Sprintf("ResourceManager object updated <%s>. Removing handler...", request.NamespacedName)))
		r.removeResourceManagerHandler(request.NamespacedName)
	}

	if resourceManager.Spec.Disabled {
		r.log.Info(trace(fmt.Sprintf("ResourceManager object disabled <%s>. Ignoring...", request.NamespacedName)))
		return ctrl.Result{}, nil
	}

	r.log.Info(trace(fmt.Sprintf("ResourceManager object added <%s>. Handler creating...", request.NamespacedName)))
	resourceManagerHandler, err := NewResourceManagerHandler(resourceManager, r.clientset, r.log)
	if err != nil {
		r.log.Error(err, fmt.Sprintf("ResourceManagerHandler object %s handler creating failed with error <%s>.", request.NamespacedName, err))
		return ctrl.Result{}, nil
	}

	// add handler to resourceManagerHandlers
	r.log.Info(trace(fmt.Sprintf("ResourceManagerHandler for <%s> registering...", request.NamespacedName)))
	r.registerAndRunResourceManagerHandler(request.NamespacedName, resourceManagerHandler)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	configLog := uzap.NewProductionEncoderConfig()
	configLog.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339Nano))
	}
	logfmtEncoder := zaplogfmt.NewEncoder(configLog)
	// Construct a new logr.logger.
	r.log = zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout), zap.Encoder(logfmtEncoder))
	log.SetLogger(r.log)

	// resourceManagerHandlers = make(map[types.NamespacedName]chan struct{})
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err.Error())
	}

	r.resourceManagerHandlers = make(map[types.NamespacedName]*ResourceManagerHandler)

	r.clientset, err = kubernetes.NewForConfig(cfg)
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
