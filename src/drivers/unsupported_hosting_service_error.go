package drivers

// UnsupportedHostingServiceError represents the error condition
// when no suitable hosting service is found in the respective Registry.
type UnsupportedHostingServiceError struct {
	registry *Registry
}
