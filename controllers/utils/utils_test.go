package utils_test

import (
	"github.com/tikalk/resource-manager/controllers/utils"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

var _ = Context("Testing utils", func() {
	Describe("testing object expiration", func() {
		It("testing IsObjExpired", func() {
			// providing now time, and 1 min expiration. expect for 60 seconds to expiration
			err, seconds := utils.IsObjExpired(time.Now(), "1m")
			Expect(err).NotTo(HaveOccurred(), "failed to calc IsIntervalOccurred")
			Expect(seconds).Should(BeNumerically("<=", 60))
		})

		It("testing IsIntervalOccurred", func() {
			// providing 15:54 time, and 15:55 expiration. expect for 60 seconds to expiration
			err, seconds := utils.IsIntervalOccurred(time.Date(2021, 8, 15, 15, 54, 0, 0, time.Local), "15:55")
			Expect(err).NotTo(HaveOccurred(), "failed to calc IsIntervalOccurred")
			Expect(seconds).To(Equal(60))
		})
	})

})

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Utils Suite",
		[]Reporter{printer.NewlineReporter{}})
}
