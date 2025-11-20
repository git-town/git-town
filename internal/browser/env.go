package browser

// name of the environment variable that defines the browser to use
// TODO: move this logic into internal/config/envconfig/load.go
const EnvVarName = "BROWSER"

// if EnvVarName content is this value, pretend that you cannot open the browser
const EnvVarNone = "(none)"
