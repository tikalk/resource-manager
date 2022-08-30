package handlers

import (
	"fmt"

	utils "github.com/tikalk/resource-manager/controllers/utils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//  keep everything in a struct

// HandledeploymentObj handle deployment objects that related to the resource-manager controller
func (d depObj) HandleDeployObj() {
	// get all the deployments with the desired selector labels
	deploy, err := GetDeployByLabel(d)
	if err != nil {
		d.L.Error(err, fmt.Sprintf("%s: cannot list deployments\n", d.Name))
	}

	if len(deploy) <= 0 {
		fmt.Printf("%s: did not found any deployments with the requested label\n", d.Name)
		return
	}

	fmt.Printf("found %d deployments with the requested label\n", len(deploy))

	for _, dep := range deploy {
		expired, secondsUntilExpire := utils.IsObjExpired(dep.CreationTimestamp, d.Spec.Condition[0].After)
		if expired {
			switch d.Spec.Action {
			case "delete":
				fmt.Printf("deployment '%s' has been expired and will be deleted \n", dep.Name)
				err := d.C.Delete(d.Ctx, dep.DeepCopy(), &client.DeleteOptions{})
				if err != nil {
					d.L.Error(err, fmt.Sprintf("cannot delete deployments\n"))
				}
				fmt.Printf("%s: deployment '%s' has been deleted \n", d.Name, dep.Name)
			}
		} else {
			fmt.Printf("%s: %d seconds has left to deployment '%s' \n", d.Name, secondsUntilExpire, dep.Name)

		}
	}
}

// GetdeploymentsByLabel get only deployments that contains a specific label
func GetDeployByLabel(d depObj) ([]appsv1.Deployment, error) {

	var listOfDeployments []appsv1.Deployment
	depListObj := &appsv1.DeploymentList{}

	if err := d.C.List(d.Ctx, depListObj, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels),
	}); err != nil {
		d.L.Error(err, fmt.Sprintf("%s: unable to fetch deployments", d.Name))
		return nil, err
	}

	for _, item := range depListObj.Items {
		listOfDeployments = append(listOfDeployments, item)
	}
	return listOfDeployments, nil
}
