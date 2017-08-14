package drivers

import "sort"

// Registry collects and manages all CodeHostingDriver instances
type Registry struct {
	drivers []*CodeHostingDriver
}

// DetermineActiveDriver determines the driver to use for the given hostname
func (r *Registry) DetermineActiveDriver(hostname string) (*CodeHostingDriver, error) {
	for _, driver := range r.drivers {
		if driver.CanBeUsed(hostname) {
			return driver, nil
		}
	}
	return nil, UnsupportedHostingServiceError{r}
}

// DriverNames returns the names of all drivers, sorted alphabetically
func (r *Registry) DriverNames() (result []string) {
	for _, driver := range r.drivers {
		result = append(result, driver.HostingServiceName)
	}
	sort.Strings(result)
	return
}

// RegisterDriver allows driver implementations to register themselves
// with the registry
func (r *Registry) RegisterDriver(driver *CodeHostingDriver) {
	r.drivers = append(r.drivers, driver)
}
