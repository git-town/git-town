package drivers

import "github.com/git-town/git-town/src/git"

// Core provides the public API for the drivers subsystem.

var registry = Registry{}

var activeDriver CodeHostingDriver

// DriverConfigurationInterface defines the driver's interface to configuration.
type DriverConfigurationInterface interface {
	GetCodeHostingDriverName() string
	GetRemoteOriginURL() string
	GetCodeHostingOriginHostname() string
}

// DriverConfiguration implements DriverConfigurationInterface.
// Exported for overrides in test.
var DriverConfiguration DriverConfigurationInterface = git.Config()

// GetActiveDriver returns the code hosting driver to use based on the git config.
func GetActiveDriver() CodeHostingDriver {
	if activeDriver == nil {
		activeDriver = GetDriver(DriverOptions{
			DriverType:     DriverConfiguration.GetCodeHostingDriverName(),
			OriginURL:      DriverConfiguration.GetRemoteOriginURL(),
			OriginHostname: DriverConfiguration.GetCodeHostingOriginHostname(),
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
