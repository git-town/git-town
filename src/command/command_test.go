package command_test

import (
	"github.com/Originate/git-town/src/command"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var res *command.Result

var _ = Describe("Run", func() {
	It("Runs the given command", func() {
		res = command.Run("echo", "foo")
		Expect(res.Output()).To(Equal("foo"))
	})
})

var _ = Describe("OutputContainsText", func() {

	BeforeEach(func() {
		res = command.Run("echo", "hello world how are you?")
	})

	It("returns true if the output contains the given text", func() {
		Expect(res.OutputContainsText("world")).To(BeTrue())
	})

	It("returns false if the output does not contain the given text", func() {
		Expect(res.OutputContainsText("zonk")).To(BeFalse())
	})
})

var _ = Describe("OutputContainsLine", func() {

	BeforeEach(func() {
		res = command.Run("echo", "hello world")
	})

	It("returns true if the output contains the given line", func() {
		Expect(res.OutputContainsLine("hello world")).To(BeTrue())
	})

	It("returns false if the output contains only parts of the given line", func() {
		Expect(res.OutputContainsLine("hello")).To(BeFalse())
	})

	It("returns false if the output does not contains the given line", func() {
		Expect(res.OutputContainsLine("zonk")).To(BeFalse())
	})
})

var _ = Describe("Err", func() {

	Context("command not found", func() {

		BeforeEach(func() {
			res = command.Run("zonk")
		})

		It("returns an error", func() {
			Expect(res.Err()).To(HaveOccurred())
		})
	})

	Context("command returns exit code", func() {

		BeforeEach(func() {
			res = command.Run("bash", "-c", "exit 2")
		})

		It("returns an error", func() {
			Expect(res.Err()).To(HaveOccurred())
		})
	})
})
