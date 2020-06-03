package drivers

import "github.com/git-town/git-town/src/git"

// Core provides the public API for the drivers subsystem.

var registry = Registry{}

var activeDriver CodeHostingDriver

// GetActiveDriver returns the code hosting driver to use based on the git config.
func GetActiveDriver() CodeHostingDriver {
	driver := TryUseGithub(git.Config())
	if driver != nil {
		return driver
	}
	if activeDriver == nil {
		activeDriver = GetDriver(DriverOptions{
			DriverType:     git.Config().GetCodeHostingDriverName(),
			OriginURL:      git.Config().GetRemoteOriginURL(),
			OriginHostname: git.Config().GetCodeHostingOriginHostname(),
		})
		if activeDriver != nil {
			activeDriver.SetAPIToken(activeDriver.GetAPIToken())
		}
	}
	return activeDriver
}

// GetDriver returns the code hosting driver to use based on given origin url.
func GetDriver(driverOptions DriverOptions) CodeHostingDriver {
	return registry.DetermineActiveDriver(driverOptions)
}

// ValidateHasDriver returns an error if there is no code hosting driver.
func ValidateHasDriver() error {
	driver := GetActiveDriver()
	if driver == nil {
		return UnsupportedHostingServiceError{&registry}
	}
	return nil
}
