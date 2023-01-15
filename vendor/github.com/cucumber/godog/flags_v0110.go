package godog

import (
	"errors"
	"flag"
	"math/rand"
	"time"

	"github.com/spf13/pflag"

	"github.com/cucumber/godog/internal/flags"
)

// Choose randomly assigns a convenient pseudo-random seed value.
// The resulting seed will be between `1-99999` for later ease of specification.
func makeRandomSeed() int64 {
	return rand.New(rand.NewSource(time.Now().UTC().UnixNano())).Int63n(99998) + 1
}

func flagSet(opt *Options) *pflag.FlagSet {
	set := pflag.NewFlagSet("godog", pflag.ExitOnError)
	flags.BindRunCmdFlags("", set, opt)
	pflag.ErrHelp = errors.New("godog: help requested")
	return set
}

// BindCommandLineFlags binds godog flags to given flag set prefixed
// by given prefix, without overriding usage
func BindCommandLineFlags(prefix string, opts *Options) {
	flagSet := pflag.CommandLine
	flags.BindRunCmdFlags(prefix, flagSet, opts)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}
