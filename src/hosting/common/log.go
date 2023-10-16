package common

// Log allows hosting adapters to print network operations to the CLI.
type Log interface {
	Start(string, ...interface{})
	Success()
	Failed(error)
}
