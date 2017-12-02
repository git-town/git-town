package command_test

import (
	"time"

	"github.com/Originate/git-town/src/command"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var cmd *command.Command

var _ = Describe("New", func() {

	BeforeEach(func() {
		cmd = command.New("cd", ".")
	})

	It("creates a new Cmd instance", func() {
		Expect(cmd).ToNot(BeNil())
	})
})

var _ = Describe("Run", func() {

	BeforeEach(func() {
		cmd = command.New("cd", ".")
	})

	It("runs the given command", func() {
		cmd.Run()
	})

	Describe("Idempotency", func() {

		BeforeEach(func() {
			cmd = command.New("ruby", "-e", "puts Time.now.to_f")
		})

		It("does not re-run the command", func() {
			cmd.Run()
			firstOutput := cmd.Output()
			time.Sleep(1 * time.Millisecond)
			cmd.Run()
			secondOutput := cmd.Output()
			Expect(firstOutput).To(Equal(secondOutput))
		})
	})
})

var _ = Describe("Output", func() {

	Describe("return value", func() {

		BeforeEach(func() {
			cmd = command.New("echo", "foo")
		})

		It("returns the output of this command", func() {
			cmd.Run()
			Expect(cmd.Output()).To(Equal("foo"))
		})

		It("runs the command if it hasn't run so far", func() {
			Expect(cmd.Output()).To(Equal("foo"))
		})
	})

	Describe("Idempotency", func() {

		BeforeEach(func() {
			cmd = command.New("ruby", "-e", "puts Time.now.to_f")
		})

		It("does not re-run the command", func() {
			firstOutput := cmd.Output()
			time.Sleep(1 * time.Millisecond)
			secondOutput := cmd.Output()
			Expect(firstOutput).To(Equal(secondOutput))
		})
	})
})

var _ = Describe("OutputContainsText", func() {

	BeforeEach(func() {
		cmd = command.New("echo", "hello world how are you?")
	})

	It("returns true if the output contains the given text", func() {
		Expect(cmd.OutputContainsText("world")).To(BeTrue())
	})

	It("returns false if the output does not contain the given text", func() {
		Expect(cmd.OutputContainsText("zonk")).To(BeFalse())
	})
})

var _ = Describe("OutputContainsLine", func() {

	BeforeEach(func() {
		cmd = command.New("echo", "hello world")
	})

	It("returns true if the output contains the given line", func() {
		Expect(cmd.OutputContainsLine("hello world")).To(BeTrue())
	})

	It("returns false if the output contains only parts of the given line", func() {
		Expect(cmd.OutputContainsLine("hello")).To(BeFalse())
	})

	It("returns false if the output does not contains the given line", func() {
		Expect(cmd.OutputContainsLine("zonk")).To(BeFalse())
	})
})

var _ = Describe("Err", func() {

	Context("command not found", func() {

		BeforeEach(func() {
			cmd = command.New("zonk")
		})

		It("returns an error", func() {
			Expect(cmd.Err()).To(HaveOccurred())
		})
	})

	Context("command returns exit code", func() {

		BeforeEach(func() {
			cmd = command.New("bash", "-c", "exit 2")
		})

		It("returns an error", func() {
			Expect(cmd.Err()).To(HaveOccurred())
		})
	})
})
