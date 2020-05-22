package drivers

import "github.com/git-town/git-town/src/git"

// Core provides the public API for the drivers subsystem

var registry = Registry{}

var activeDriver CodeHostingDriver

var gitConfig = git.Config()

// GetActiveDriver returns the code hosting driver to use based on the git config
func GetActiveDriver() CodeHostingDriver {
	var originHostname string
	originURL := gitConfig.GetRemoteOriginURL()
	if gitConfig.GetCodeHostingOriginHostname() != "" {
		originHostname = gitConfig.GetCodeHostingOriginHostname()
	} else {
		originHostname = gitConfig.GetURLHostname(originURL)
	}
	if activeDriver == nil {
		activeDriver = GetDriver(DriverOptions{
			DriverType:     gitConfig.GetCodeHostingDriverName(),
			OriginURL:      originURL,
			OriginHostname: originHostname,
		})
	}
	return activeDriver
}

// GetDriver returns the code hosting driver to use based on given origin url
func GetDriver(driverOptions DriverOptions) CodeHostingDriver {
	return registry.DetermineActiveDriver(driverOptions)
}

// ValidateHasDriver returns an error if there is no code hosting driver
func ValidateHasDriver() error {
	driver := GetActiveDriver()
	if driver == nil {
		return UnsupportedHostingServiceError{&registry}
	}
	return nil
}
