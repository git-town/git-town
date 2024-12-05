package flags

import (
	"github.com/spf13/pflag"
)

// BindRunCmdFlags is an internal func to bind run subcommand flags.
func BindRunCmdFlags(prefix string, flagSet *pflag.FlagSet, opts *Options) {
	if opts.Concurrency == 0 {
		opts.Concurrency = 1
	}

	if opts.Format == "" {
		opts.Format = "pretty"
	}

	flagSet.BoolVar(&opts.NoColors, prefix+"no-colors", opts.NoColors, "disable ansi colors")
	flagSet.IntVarP(&opts.Concurrency, prefix+"concurrency", "c", opts.Concurrency, "run the test suite with concurrency")
	flagSet.StringVarP(&opts.Tags, prefix+"tags", "t", opts.Tags, `filter scenarios by tags, expression can be:
  "@wip"           run all scenarios with wip tag
  "~@wip"          exclude all scenarios with wip tag
  "@wip && ~@new"  run wip scenarios, but exclude new
  "@wip,@undone"   run wip or undone scenarios`)
	flagSet.StringVarP(&opts.Format, prefix+"format", "f", opts.Format, `will write a report according to the selected formatter

usage:
  -f <formatter>
  will use the formatter and write the report on stdout
  -f <formatter>:<file_path>
  will use the formatter and write the report to the file path

built-in formatters are:
  progress  prints a character per step
  cucumber  produces a Cucumber JSON report
  events    produces JSON event stream, based on spec: 0.1.0
  junit     produces JUnit compatible XML report
  pretty    prints every feature with runtime statuses
 `)

	flagSet.BoolVarP(&opts.ShowStepDefinitions, prefix+"definitions", "d", opts.ShowStepDefinitions, "print all available step definitions")
	flagSet.BoolVar(&opts.StopOnFailure, prefix+"stop-on-failure", opts.StopOnFailure, "stop processing on first failed scenario")
	flagSet.BoolVar(&opts.Strict, prefix+"strict", opts.Strict, "fail suite when there are pending or undefined steps")

	flagSet.Int64Var(&opts.Randomize, prefix+"random", opts.Randomize, `randomly shuffle the scenario execution order
  --random
specify SEED to reproduce the shuffling from a previous run
  --random=5738`)
	flagSet.Lookup(prefix + "random").NoOptDefVal = "-1"
}
