package drivers

import "github.com/git-town/git-town/src/git"

// Core provides the public API for the drivers subsystem.

var registry = Registry{}

var activeDriver CodeHostingDriver

// ConfigurationInterface defines the drivers's interface to configuration.
type ConfigurationInterface interface {
	GetCodeHostingDriverName() string
	GetRemoteOriginURL() string
	GetCodeHostingOriginHostname() string
	GetURLHostname(string) string
	GetURLRepositoryName(string) string
}

// Configuration implements ConfigurationInterface.
// Exported for overrides in test.
var Configuration ConfigurationInterface = git.Config()

// GetActiveDriver returns the code hosting driver to use based on the git config.
func GetActiveDriver() CodeHostingDriver {
	if activeDriver == nil {
		activeDriver = GetDriver(DriverOptions{
			DriverType:     Configuration.GetCodeHostingDriverName(),
			OriginURL:      Configuration.GetRemoteOriginURL(),
			OriginHostname: Configuration.GetCodeHostingOriginHostname(),
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
