package drivers

import "fmt"

// UnsupportedHostingServiceError represents the error condition
// when no suitable hosting service is found in the respective Registry
type UnsupportedHostingServiceError struct {
	registry *Registry
}

func (e UnsupportedHostingServiceError) Error() string {
	result := "Unsupported hosting service\n\nThis command requires hosting on one of these services:\n"
	for _, driverName := range e.registry.DriverNames() {
		result = fmt.Sprintf("%s* %s\n", result, driverName)
	}
	return result
}
