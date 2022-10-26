package utils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

var _ = Context("Testing utils", func() {
	Describe("testing object expiration", func() {
		It("testing IsObjExpired", func() {
			// creationTime := time.Parse("", )
			// err, seconds := utils.IsObjExpired()
			// Expect(err).NotTo(HaveOccurred(), "failed to calc IsIntervalOccurred")
			// Expect(seconds).To(Equal(60))

			// Expect().To(Equal("delete"))
			// Expect(rmObj.Spec.Condition.Timeframe).To(Equal("1s"))
		})

		It("testing IsIntervalOccurred", func() {

			// err, seconds := utils.IsIntervalOccurred("16:54")
			// Expect(err).NotTo(HaveOccurred(), "failed to calc IsIntervalOccurred")
			// Expect(seconds).To(Equal(60))

			// Expect().To(Equal("delete"))
			// Expect(rmObj.Spec.Condition.Timeframe).To(Equal("1s"))

		})
	})
})

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Utils Suite",
		[]Reporter{printer.NewlineReporter{}})
}
