package command_test

import (
	"github.com/Originate/git-town/src/command"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("New", func() {

	Context("with a single parameter", func() {
		It("creates a new Cmd instance with the given command", func() {
			cmd := command.New("ls")
			Expect(cmd.Name).To(Equal("ls"))
			Expect(cmd.Args).To(Equal([]string{}))
		})
	})

	Context("with multiple parameters", func() {
		It("creates a new Cmd instance with the given command and arguments", func() {
			cmd := command.New("ls", "-la", "*.go")
			Expect(cmd.Name).To(Equal("ls"))
			Expect(cmd.Args).To(Equal([]string{"-la", "*.go"}))
		})
	})

})

var _ = Describe("Output", func() {

	It("returns the output of this command", func() {
		cmd := command.New("echo", "foo")
		Expect(cmd.Output()).To(Equal("foo"))
	})
})

var _ = Describe("OutputContainsText", func() {

	It("returns whether the output contains the given text", func() {
		cmd := command.New("echo", "hello world how are you?")
		Expect(cmd.OutputContainsText("world")).To(BeTrue())
	})
})
