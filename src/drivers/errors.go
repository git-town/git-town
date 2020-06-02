package drivers

import "errors"

var ErrNotSupported = errors.New("not supported")
var ErrNoPullRequestFound = errors.New("no pull request found")
