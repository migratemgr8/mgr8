package global

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Check Command", func() {
	var (
		subject Database
	)

	Context("FromStr", func() {
		When("not used", func() {
			It("uses default value", func() {
				Expect(subject).To(Equal(Postgres))
			})
		})

		for k, v := range toStr {
			When(fmt.Sprintf("str is %s", v), func() {
				subject.FromStr(v)
				It(fmt.Sprintf("updates to %s domain", v), func() {
					Expect(subject).To(Equal(k))
				})
			})
		}
	})

	Context("ToStr", func() {
		for k, v := range toStr {
			When(v, func() {
				subject = k
				Expect(fmt.Sprintf("%s", subject)).To(Equal(v))
			})
		}
	})
})
