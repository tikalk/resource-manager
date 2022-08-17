package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	resourcemanagmentv1alpha1 "github.com/tikalk/resource-manager/api/v1alpha1"
)

var _ = Context("Inside of a ResourceManager", func() {
	ctx := context.TODO()
	SetupTest(ctx)

	Describe("when no existing resources exist", func() {
		It("should create a new expiry resource manager", func() {
			myResourceManagerObj := &resourcemanagmentv1alpha1.ResourceManager{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-resource-manager",
					Namespace: "default",
				},
				Spec: resourcemanagmentv1alpha1.ResourceManagerSpec{
					Active:    true,
					DryRun:    false,
					Resources: "namespace",
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"name": "managed-namespace",
						},
					},
					Action: "delete",
					Condition: []resourcemanagmentv1alpha1.ExpiryCondition{{
						Condition: resourcemanagmentv1alpha1.Condition{
							Type: "expiry",
						},
						After: "1m",
					}},
				},
			}

			err := k8sClient.Create(ctx, myResourceManagerObj)
			Expect(err).NotTo(HaveOccurred(), "failed to create test ResourceManager resource")

			rmObj := &resourcemanagmentv1alpha1.ResourceManager{}
			Eventually(
				getResourceFunc(ctx, client.ObjectKey{Name: "test-resource-manager", Namespace: myResourceManagerObj.Namespace}, rmObj),
				time.Second*5, time.Millisecond*500).Should(BeNil())

			Expect(rmObj.Spec.Action).To(Equal("delete"))
			Expect(rmObj.Spec.Condition[0].After).To(Equal("1m"))
			// Expect(rmObj.Spec.Condition[0].After).To(Equal("2m"))
		})
	})
})

func getResourceFunc(ctx context.Context, key client.ObjectKey, obj client.Object) func() error {
	return func() error {
		return k8sClient.Get(ctx, key, obj)
	}
}
