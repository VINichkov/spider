package block_list_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBlockList(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BlockList")
}

var _ = Describe("BlockList", func() {
	list := NewBlockList()
	LocationOne := 5
	LocationTwo := 6
	list.Push(LocationOne)
	list.Push(LocationTwo)

	AfterSuite(func() {
		list.Close()
	})

	It("should include two locations", func() {
		resultOne := list.Include(LocationTwo)
		resultTwo := list.Include(LocationOne)

		Expect(resultOne).To(BeTrue())
		Expect(resultTwo).To(BeTrue())
	})

	It("should not include a location", func() {
		Expect(list.Include(9)).To(BeFalse())
	})

})
