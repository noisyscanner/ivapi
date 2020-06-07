package tokens

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokens", func() {
	var _ = Describe("GenerateToken", func() {
		It("should return a token 16 chars in length", func() {
			Expect(GenerateToken()).To(HaveLen(16))
		})
	})
})
