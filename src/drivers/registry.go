package drivers

import (
	"fmt"
	"os"
	"sort"
)

// Registry is the micro-co
type Registry struct {
	drivers []*CodeHostingDriver
}

// RegisterDriver allows driver implementations to register themselves
// with the registry
func (c *Registry) RegisterDriver(driver *CodeHostingDriver) {
	c.drivers = append(c.drivers, driver)
}

// DetermineActiveDriver determines the driver to use for the current repository
func (c *Registry) DetermineActiveDriver(hostname string) *CodeHostingDriver {
	for _, driver := range c.drivers {
		if driver.CanBeUsed(hostname) {
			return driver
		}
	}

	fmt.Println("Unsupported hosting service")
	fmt.Println()
	fmt.Println("This command requires hosting on one of these services:")
	driverNames := []string{}
	for _, driver := range c.drivers {
		driverNames = append(driverNames, driver.HostingServiceName)
	}
	sort.Strings(driverNames)
	for _, driverName := range driverNames {
		fmt.Printf("* %s\n", driverName)
	}
	os.Exit(1)
	return nil
}
