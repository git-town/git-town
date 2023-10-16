package common

// Log allows hosting adapters to print network operations to the CLI.
// TODO: target the respective struct in the CLI package directly?
type Log interface {
	Start(string, ...interface{})
	Success()
	Failed(error)
}
