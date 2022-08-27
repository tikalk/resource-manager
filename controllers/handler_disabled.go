package controllers

import (
	"context"
	"fmt"
	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

//+kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=*,resources=namespaces/finalizers,verbs=update

func HandleDisabled(ctx context.Context,
	stopper chan struct{},
	action string,
	condition []resourcemanagmentv1alpha1.ExpiryCondition,
	managedResource resourcemanagmentv1alpha1.ResourceSelector,
	namespace string,
	clientset *kubernetes.Clientset) {

	l := log.FromContext(ctx)
	l.Info(fmt.Sprintf("HandleDisabled begin"))
	<-stopper
	l.Info(fmt.Sprintf("HandleDisabled end"))

}
