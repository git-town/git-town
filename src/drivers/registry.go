package drivers

import (
	"sort"
)

// Registry collects and manages all CodeHostingDriver instances
type Registry struct {
	drivers []CodeHostingDriver
}

// DetermineActiveDriver determines the driver to use for the given hostname
func (r *Registry) DetermineActiveDriver(driverOptions DriverOptions) CodeHostingDriver {
	for _, driver := range r.drivers {
		driver.SetOriginURL(driverOptions.OriginURL)
		if driverOptions.OriginHostname != "" {
			driver.SetOriginHostname(driverOptions.OriginHostname)
		}
		if driver.CanBeUsed(driverOptions.DriverType) {
			return driver
		}
	}
	return nil
}

// DriverNames returns the names of all drivers, sorted alphabetically
func (r *Registry) DriverNames() (result []string) {
	for _, driver := range r.drivers {
		result = append(result, driver.HostingServiceName())
	}
	sort.Strings(result)
	return
}

// RegisterDriver allows driver implementations to register themselves
// with the registry
func (r *Registry) RegisterDriver(driver CodeHostingDriver) {
	r.drivers = append(r.drivers, driver)
}
