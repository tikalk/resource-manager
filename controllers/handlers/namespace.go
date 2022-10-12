package handlers

import (
	"fmt"
	"time"

	"github.com/tikalk/resource-manager/controllers/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HandleNamespaceObj handle namespace objects that related to the resource-manager controller
func (o Obj) HandleNamespaceObj() {
	// get all the namespaces with the desired selector labels
	namespacesToHandle, err := o.GetNamespacesByLabel()
	if err != nil {
		o.l.Error(err, fmt.Sprintf("%s: cannot list namespaces\n", o.Name))
		return
	}

	if len(namespacesToHandle) <= 0 {
		fmt.Printf("%s: did not found any namespace with the requested label\n", o.Name)
		return
	}

	for _, ns := range namespacesToHandle {
		switch o.rm.Spec.Condition[0].Type {
		case "expiry":
			o.handleExpiry(ns)
		case "timeframe":
			o.handleTimeframe(ns)
		}

	}
	time.Sleep(10 * time.Second)
}

// GetNamespacesByLabel get only namespaces that contains a specific label
func (o Obj) GetNamespacesByLabel() ([]v1.Namespace, error) {

	var listOfNamespaces []v1.Namespace
	nsListObj := &v1.NamespaceList{}

	if err := o.c.List(o.ctx, nsListObj, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(o.rm.Spec.Selector.MatchLabels),
	}); err != nil {
		o.l.Error(err, fmt.Sprintf("%s: unable to fetch namespaces", o.Name))
		return nil, err
	}

	for _, item := range nsListObj.Items {
		listOfNamespaces = append(listOfNamespaces, item)
	}
	return listOfNamespaces, nil
}

// deleteNamespace delete namespace obj
func (o Obj) deleteNamespace(namespace v1.Namespace) {
	err := o.c.Delete(o.ctx, namespace.DeepCopy(), &client.DeleteOptions{})
	if err != nil {
		o.l.Error(err, fmt.Sprintf("cannot delete namespaces\n"))
	}
	time.Sleep(5 * time.Second)
	fmt.Printf("%s: namespace '%s' has been deleted \n", o.Name, namespace.Name)

}

// handleTimeframe handle timeframe type
func (o Obj) handleTimeframe(namespace v1.Namespace) {
	fmt.Printf("namespace '%s' will be deleted at timeframe: %s  \n", namespace.Name, o.rm.Spec.Condition[0].Timeframe)
	err, doesIntervalOccurred := utils.IsIntervalOccurred(o.rm.Spec.Condition[0].Timeframe)
	if err != nil {
		o.l.Error(err, fmt.Sprintf("cannot calculate timeframe\n"))
		return
	}
	if doesIntervalOccurred {
		switch o.rm.Spec.Action {
		case "delete":
			fmt.Printf("namespace '%s' is in timeframe and will be deleted \n", namespace.Name)
			o.deleteNamespace(namespace)
		}
	}
}

// handleExpiry handle expiry type
func (o Obj) handleExpiry(namespace v1.Namespace) {
	expired, secondsUntilExpire := utils.IsObjExpired(namespace.CreationTimestamp, o.rm.Spec.Condition[0].After)
	if expired {
		switch o.rm.Spec.Action {
		case "delete":
			fmt.Printf("namespace '%s' has been expired and will be deleted \n", namespace.Name)
			o.deleteNamespace(namespace)
			fmt.Printf("%s: namespace '%s' has been deleted \n", o.Name, namespace.Name)
		}
	} else {
		fmt.Printf("%s: %d seconds has left to namespace '%s' \n", o.Name, secondsUntilExpire, namespace.Name)
	}
}
