package handlers

import (
	"fmt"
	utils "github.com/tikalk/resource-manager/controllers/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//  keep everything in a struct

// HandleNamespaceObj handle namespace objects that related to the resource-manager controller
func (o Obj) HandleNamespaceObj() {
	// get all the namespaces with the desired selector labels
	namespaces, err := GetNamespacesByLabel(o)
	if err != nil {
		o.L.Error(err, fmt.Sprintf("%s: cannot list namespaces\n", o.Name))
	}

	if len(namespaces) <= 0 {
		fmt.Printf("%s: did not found any namespaces with the requested label\n", o.Name)
		return
	}

	fmt.Printf("found %d namespaces with the requested label\n", len(namespaces))

	for _, namespace := range namespaces {
		expired, secondsUntilExpire := utils.IsObjExpired(namespace.CreationTimestamp, o.Spec.Condition[0].After)
		if expired {
			switch o.Spec.Action {
			case "delete":
				fmt.Printf("namespace '%s' has been expired and will be deleted \n", namespace.Name)
				err := o.C.Delete(o.Ctx, namespace.DeepCopy(), &client.DeleteOptions{})
				if err != nil {
					o.L.Error(err, fmt.Sprintf("cannot delete namespaces\n"))
				}
				fmt.Printf("%s: namespace '%s' has been deleted \n", o.Name, o.Name)
			}
		} else {
			fmt.Printf("%s: %d seconds has left to namespace '%s' \n", o.Name, secondsUntilExpire, namespace.Name)

		}
	}
}

// GetNamespacesByLabel get only namespaces that contains a specific label
func GetNamespacesByLabel(o Obj) ([]v1.Namespace, error) {

	var listOfNamespaces []v1.Namespace
	nsListObj := &v1.NamespaceList{}

	if err := o.C.List(o.Ctx, nsListObj, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(o.Spec.Selector.MatchLabels),
	}); err != nil {
		o.L.Error(err, fmt.Sprintf("%s: unable to fetch namespaces", o.Name))
		return nil, err
	}

	for _, item := range nsListObj.Items {
		listOfNamespaces = append(listOfNamespaces, item)
	}
	return listOfNamespaces, nil
}
