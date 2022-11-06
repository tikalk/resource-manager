package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
)

var _ = Context("Inside of a ResourceManager", func() {
	ctx := context.TODO()
	SetupTest(ctx)

	Describe("when no existing resources exist", func() {
		It("when creating a new resource manager object and a namespace obj and wait 20 seconds", func() {

			// create new namespace object
			myNamespaceObj := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
				Name: "xxx",
				Labels: map[string]string{
					"env": "development",
				},
			},
			}

			err := k8sClient.Create(ctx, myNamespaceObj)
			Expect(err).NotTo(HaveOccurred(), "failed to create test 'namespace' resource")

			//validate creation
			nsObj, _ := getNsByName(ctx, myNamespaceObj.Name)
			Expect(string(nsObj.Status.Phase)).Should(Equal("Active"))

			myResourceManagerObj := &resourcemanagmentv1alpha1.ResourceManager{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "resource-manager-sample",
					Namespace: "xxx",
				},
				Spec: resourcemanagmentv1alpha1.ResourceManagerSpec{
					Disabled:     false,
					DryRun:       false,
					ResourceKind: "Namespace",
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"env": "development",
						},
					},
					Action: "delete",
					Condition: resourcemanagmentv1alpha1.Expiration{
						ExpireAfter: "5s",
					},
				},
			}

			err = k8sClient.Create(ctx, myResourceManagerObj)
			Expect(err).NotTo(HaveOccurred(), "failed to create test 'ResourceManager' resource")

			rmObj := &resourcemanagmentv1alpha1.ResourceManager{}
			Eventually(
				getResourceFunc(ctx, client.ObjectKey{Name: "resource-manager-sample",
					Namespace: myResourceManagerObj.Namespace}, rmObj),
				time.Second*5, time.Millisecond*500).Should(BeNil())

			Expect(rmObj.Spec.Action).To(Equal("delete"))
			Expect(rmObj.Spec.Condition.ExpireAfter).To(Equal("5s"))
			time.Sleep(10 * time.Second)

			//nsListObj := &v1.NamespaceList{}
			//
			//if err := k8sClient.List(ctx, nsListObj, &client.ListOptions{
			//	LabelSelector: labels.SelectorFromSet(map[string]string{
			//		"env": "development",
			//	})}); err != nil {
			//	fmt.Print("unable to fetch namespaces")
			//}
			//
			//if len(nsListObj.Items) != 0 {
			//	for _, item := range nsListObj.Items {
			//		fmt.Printf("\n found namespace: %s lables: %s \n", item.Name, item.Labels)
			//	}
			//} else {
			//	fmt.Println("did not find namespaces")
			//}

			//validate deletion
			//err = k8sClient.Delete(ctx, myNamespaceObj)
			//Expect(err).NotTo(HaveOccurred(), "failed to manually delete test 'xxx' resource")

			ns, _ := getNsByName(ctx, "xxx")
			Expect(string(ns.Status.Phase)).To(Not(Equal("Active")))

		})

	})
})

func getResourceFunc(ctx context.Context, key client.ObjectKey, obj client.Object) func() error {
	return func() error {
		return k8sClient.Get(ctx, key, obj)
	}
}

func getNsByName(ctx context.Context, name string) (*v1.Namespace, error) {
	nsObj := &v1.Namespace{}
	if err := k8sClient.Get(ctx, client.ObjectKey{Name: name}, nsObj); err != nil {
		return nil, nil
	}
	return nsObj, nil
}
