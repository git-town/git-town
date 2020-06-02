package cmd

import "errors"

var ErrNoFeatureBranch = errors.New("no feature branch")
var ErrBranchMissing = errors.New("branch missing")
var ErrInvalidValue = errors.New("invalid value")
