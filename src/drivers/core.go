package drivers

import "github.com/git-town/git-town/src/git"

// Core provides the public API for the drivers subsystem

var registry = Registry{}

var activeDriver CodeHostingDriver

var GitConfig = git.Config()

// GetActiveDriver returns the code hosting driver to use based on the git config
func GetActiveDriver() CodeHostingDriver {
	if activeDriver == nil {
		activeDriver = GetDriver()
	}
	return activeDriver
}

// GetDriver returns the code hosting driver to use based on given origin url
func GetDriver() CodeHostingDriver {
	var originHostname string
	originURL := GitConfig.GetRemoteOriginURL()
	if GitConfig.GetCodeHostingOriginHostname() != "" {
		originHostname = GitConfig.GetCodeHostingOriginHostname()
	} else {
		originHostname = GetURLHostname(originURL)
	}
	return registry.DetermineActiveDriver(DriverOptions{
		DriverType:     GitConfig.GetCodeHostingDriverName(),
		OriginURL:      originURL,
		OriginHostname: originHostname,
	})
}

// ValidateHasDriver returns an error if there is no code hosting driver
func ValidateHasDriver() error {
	driver := GetActiveDriver()
	if driver == nil {
		return UnsupportedHostingServiceError{&registry}
	}
	return nil
}
