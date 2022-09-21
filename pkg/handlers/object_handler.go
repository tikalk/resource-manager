package handlers

import (
	v1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"
)

type ObjectHandler struct {
	resourceManager *v1alpha1.ResourceManager
	object          interface{}

	stopper chan struct{}
}

func NewObjectHandler(res *v1alpha1.ResourceManager, obj interface{}) (*ObjectHandler, types.NamespacedName, error) {
	// extract the NamespacedName of the object for storage
	fullname := extractFullname(res.Spec.ResourceKind, obj)

	// return the object handler
	o := &ObjectHandler{
		object:          obj,
		resourceManager: res,
	}
	return o, fullname, nil
}

func extractFullname(kind string, obj interface{}) (fullname types.NamespacedName) {
	switch kind {
	case "Namespace":
		fullname = types.NamespacedName{Name: obj.(*v1.Namespace).Name, Namespace: obj.(*v1.Namespace).Namespace}
	case "Deployment":
		fullname = types.NamespacedName{Name: obj.(*appsv1.Deployment).Name, Namespace: obj.(*appsv1.Deployment).Namespace}
	default:
		// TODO: print error
	}
	return fullname
}

func extractPhase(kind string, obj interface{}) (fullname types.NamespacedName) {
	switch kind {
	case "Namespace":
		fullname = types.NamespacedName{Name: obj.(*v1.Namespace).Name, Namespace: obj.(*v1.Namespace).Namespace}
	case "Deployment":
		fullname = types.NamespacedName{Name: obj.(*appsv1.Deployment).Name, Namespace: obj.(*appsv1.Deployment).Namespace}
	default:
		// TODO: print error
	}
	return fullname
}

// func (o *ObjectHandler) getObjType() string {
// 	return o.resourceManager.Spec.ResourceKind
// }

func (o *ObjectHandler) Run() {

}

func (o *ObjectHandler) Stop() {
	close(o.stopper)
}
