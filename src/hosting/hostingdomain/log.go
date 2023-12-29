package hostingdomain

// Log allows hosting adapters to print network operations to the CLI.
type Log interface {
	Start(template string, data ...interface{})
	Success()
	Failed(err error)
}
