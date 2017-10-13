package drivers_test

import (
	"github.com/Originate/git-town/src/drivers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CodeHostingDriver", func() {
	Describe("GetDriver", func() {
		Describe("with driver override", func() {
			It("works with bitbucket", func() {
				result := drivers.GetDriver(drivers.DriverOptions{
					DriverType: "bitbucket",
					OriginURL:  "git@self-hosted-bitbucket.com:Originate/git-town.git",
				})
				Expect(result).NotTo(BeNil())
				Expect(result.HostingServiceName()).To(Equal("Bitbucket"))
				Expect(result.GetRepositoryURL()).To(Equal("https://self-hosted-bitbucket.com/Originate/git-town"))
			})

			It("works with github", func() {
				result := drivers.GetDriver(drivers.DriverOptions{
					DriverType: "github",
					OriginURL:  "git@self-hosted-github.com:Originate/git-town.git",
				})
				Expect(result).NotTo(BeNil())
				Expect(result.HostingServiceName()).To(Equal("Github"))
				Expect(result.GetRepositoryURL()).To(Equal("https://self-hosted-github.com/Originate/git-town"))
			})

			It("works with gitlab", func() {
				result := drivers.GetDriver(drivers.DriverOptions{
					DriverType: "gitlab",
					OriginURL:  "git@self-hosted-gitlab.com:Originate/git-town.git",
				})
				Expect(result).NotTo(BeNil())
				Expect(result.HostingServiceName()).To(Equal("Gitlab"))
				Expect(result.GetRepositoryURL()).To(Equal("https://self-hosted-gitlab.com/Originate/git-town"))
			})
		})

		Describe("with origin hostname override", func() {
			It("works with bitbucket", func() {
				result := drivers.GetDriver(drivers.DriverOptions{
					OriginURL:      "git@my-ssh-identity.com:Originate/git-town.git",
					OriginHostname: "bitbucket.org",
				})
				Expect(result).NotTo(BeNil())
				Expect(result.HostingServiceName()).To(Equal("Bitbucket"))
				Expect(result.GetRepositoryURL()).To(Equal("https://bitbucket.org/Originate/git-town"))
			})

			It("works with github", func() {
				result := drivers.GetDriver(drivers.DriverOptions{
					OriginURL:      "git@my-ssh-identity.com:Originate/git-town.git",
					OriginHostname: "github.com",
				})
				Expect(result).NotTo(BeNil())
				Expect(result.HostingServiceName()).To(Equal("Github"))
				Expect(result.GetRepositoryURL()).To(Equal("https://github.com/Originate/git-town"))
			})

			It("works with gitlab", func() {
				result := drivers.GetDriver(drivers.DriverOptions{
					OriginURL:      "git@my-ssh-identity.com:Originate/git-town.git",
					OriginHostname: "gitlab.com",
				})
				Expect(result).NotTo(BeNil())
				Expect(result.HostingServiceName()).To(Equal("Gitlab"))
				Expect(result.GetRepositoryURL()).To(Equal("https://gitlab.com/Originate/git-town"))
			})
		})
	})
})
