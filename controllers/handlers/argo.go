package handlers

//
//import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//
////
//import (
//	"fmt"
//	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
//	"time"
//)
//
//// HandleNamespaceObj handle namespace objects that related to the resource-manager controller
//func (o ResourceManagerHandler) HandleArgoObj() {
//	// get all the namespaces with the desired selector labels
//	argocdObjToHandle, err := o.GetArgoAppl()
//
//	if err != nil {
//		o.L.Error(err, fmt.Sprintf("%s: cannot list namespaces\n", o.Name))
//		return
//	}
//
//	if len(argocdObjToHandle) <= 0 {
//		fmt.Printf("%s: did not found any namespace with the requested label\n", o.Name)
//		return
//	}
//
//	fmt.Println(argocdObjToHandle)
//	time.Sleep(5 * time.Second)
//}
//
//func (o ResourceManagerHandler) ListArgocdApplications(namespace string) ([]argocdv1alpha1.Application, error) {
//	appList, err := o.C.Argocd.ArgoprojV1alpha1().Applications(namespace).
//		List(context.TODO(), metav1.ListOptions{})
//	if err != nil {
//		return nil, err
//	}
//	return appList.Items, nil
//}
//
////
////func (o ResourceManagerHandler) GetArgoAppl() ([]argocdv1alpha1.Application, error) {
////	var listOfArgoObj []argocdv1alpha1.Application
////	list := &argocdv1alpha1.ApplicationList{}
////	if err := o.C.List(o.Ctx, list, &client.ListOptions{Namespace: "argocd"}); err != nil {
////		o.L.Error(err, fmt.Sprintf("%s: unable to fetch namespaces", o.Name))
////		return nil, err
////	}
////
////	for _, item := range list.Items {
////		listOfArgoObj = append(listOfArgoObj, item)
////	}
////	return listOfArgoObj, nil
////}
