# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org).

This document is formatted according to the principles of [Keep A CHANGELOG](http://keepachangelog.com).

## Unreleased

## [v0.15.0]

### Added
- Improved the type checking of step return types and improved the error messages - ([647](https://github.com/cucumber/godog/pull/647) - [johnlon](https://github.com/johnlon))
- Ambiguous step definitions will now be detected when strict mode is activated - ([636](https://github.com/cucumber/godog/pull/636)/([648](https://github.com/cucumber/godog/pull/648) - [johnlon](https://github.com/johnlon))
- Provide support for attachments / embeddings including a new example in the examples dir - ([623](https://github.com/cucumber/godog/pull/623) - [johnlon](https://github.com/johnlon))

### Changed
- Formatters now have a `Close` method and associated `io.Writer` changed to `io.WriteCloser`.

## [v0.14.1]

### Added
- Provide testing.T-compatible interface on test context, allowing usage of assertion libraries such as testify's assert/require - ([571](https://github.com/cucumber/godog/pull/571) - [mrsheepuk](https://github.com/mrsheepuk))
- Created releasing guidelines - ([608](https://github.com/cucumber/godog/pull/608) - [glibas](https://github.com/glibas))

### Fixed
- Step duration calculation - ([616](https://github.com/cucumber/godog/pull/616) - [iaroslav-ciupin](https://github.com/iaroslav-ciupin))
- Invalid memory address or nil pointer dereference in RetrieveFeatures - ([566](https://github.com/cucumber/godog/pull/566) - [corneldamian](https://github.com/corneldamian))

## [v0.14.0]
### Added
- Improve ErrSkip handling, add test for Summary and operations order ([584](https://github.com/cucumber/godog/pull/584) - [vearutop](https://github.com/vearutop))

### Fixed
- Remove line overwriting for scenario outlines in cucumber formatter ([605](https://github.com/cucumber/godog/pull/605) - [glibas](https://github.com/glibas))
- Remove duplicate warning message ([590](https://github.com/cucumber/godog/pull/590) - [vearutop](https://github.com/vearutop))
- updated base formatter to set a scenario as passed unless there exist ([582](https://github.com/cucumber/godog/pull/582) - [roskee](https://github.com/roskee))

### Changed
- Update test.yml ([583](https://github.com/cucumber/godog/pull/583) - [vearutop](https://github.com/vearutop))
  
## [v0.13.0]
### Added
- Support for reading feature files from an `fs.FS` ([550](https://github.com/cucumber/godog/pull/550) - [tigh-latte](https://github.com/tigh-latte))
- Added keyword functions. ([509](https://github.com/cucumber/godog/pull/509) - [otrava7](https://github.com/otrava7))
- Prefer go test to use of godog cli in README ([548](https://github.com/cucumber/godog/pull/548) - [danielhelfand](https://github.com/danielhelfand))
- Use `fs.FS` abstraction for filesystem ([550](https://github.com/cucumber/godog/pull/550) - [tigh-latte](https://github.com/tigh-latte))
- Cancel context for each scenario ([514](https://github.com/cucumber/godog/pull/514) - [draganm](https://github.com/draganm))

### Fixed
- Improve hooks invocation flow ([568](https://github.com/cucumber/godog/pull/568) - [vearutop](https://github.com/vearutop))
- Result of testing.T respect strict option ([539](https://github.com/cucumber/godog/pull/539) - [eiel](https://github.com/eiel))

### Changed
- BREAKING CHANGE, upgraded cucumber and messages dependencies = ([515](https://github.com/cucumber/godog/pull/515) - [otrava7](https://github.com/otrava7))

## [v0.12.6]
### Changed
- Each scenario is run with a cancellable `context.Context` which is cancelled at the end of the scenario. ([514](https://github.com/cucumber/godog/pull/514) - [draganm](https://github.com/draganm))
- README example is updated with `context.Context` and `go test` usage. ([477](https://github.com/cucumber/godog/pull/477) - [vearutop](https://github.com/vearutop))
- Removed deprecation of `godog.BindFlags`. ([498](https://github.com/cucumber/godog/pull/498) - [vearutop](https://github.com/vearutop))
- Pretty Print when using rules. ([480](https://github.com/cucumber/godog/pull/480) - [dumpsterfireproject](https://github.com/dumpsterfireproject))

### Fixed
- Fixed a bug which would ignore the context returned from a substep.([488](https://github.com/cucumber/godog/pull/488) - [wichert](https://github.com/wichert))
- Fixed a bug which would cause a panic when using the pretty formatter with a feature that contained a rule. ([480](https://github.com/cucumber/godog/pull/480) - [dumpsterfireproject](https://github.com/dumpsterfireproject))
- Multiple invocations of AfterScenario hooks in case of undefined steps. ([494](https://github.com/cucumber/godog/pull/494) - [vearutop](https://github.com/vearutop))
- Add a check for missing test files and raise a more helpful error. ([468](https://github.com/cucumber/godog/pull/468) - [ALCooper12](https://github.com/ALCooper12))
- Fix version subcommand. Do not print usage if run subcommand fails. ([475](https://github.com/cucumber/godog/pull/475) - [coopernurse](https://github.com/coopernurse))

### Added
- Add new option for created features with parsing from byte slices. ([476](https://github.com/cucumber/godog/pull/476) - [akaswenwilk](https://github.com/akaswenwilk))

### Deprecated
- `godog` CLI tool prints deprecation warning. ([489](https://github.com/cucumber/godog/pull/489) - [vearutop](https://github.com/vearutop))

## [v0.12.5]
### Changed
- Changed underlying cobra command setup to return errors instead of calling `os.Exit` directly to enable simpler testing. ([454](https://github.com/cucumber/godog/pull/454) - [mxygem](https://github.com/mxygem))
- Remove use of deprecated methods from `_examples`. ([460](https://github.com/cucumber/godog/pull/460) - [ricardogarfe](https://github.com/ricardogarfe))

### Fixed
- Support for go1.18 in `godog` cli mode ([466](https://github.com/cucumber/godog/pull/466) - [vearutop](https://github.com/vearutop))

## [v0.12.4]
### Added
- Allow suite-level configuration of steps and hooks ([453](https://github.com/cucumber/godog/pull/453) - [vearutop](https://github.com/vearutop))

## [v0.12.3]
### Added
- Automated binary releases with GitHub Actions ([437](https://github.com/cucumber/godog/pull/437) - [vearutop](https://github.com/vearutop))
- Automated binary versioning with `go install` ([437](https://github.com/cucumber/godog/pull/437) - [vearutop](https://github.com/vearutop))
- Module with local replace in examples ([437](https://github.com/cucumber/godog/pull/437) - [vearutop](https://github.com/vearutop))

### Changed
- suggest to use `go install` instead of the deprecated `go get` to install the `godog` binary ([449](https://github.com/cucumber/godog/pull/449) - [dmitris](https://github.com/dmitris))

### Fixed
- After Scenario hook is called before After Step ([444](https://github.com/cucumber/godog/pull/444) - [vearutop](https://github.com/vearutop))
- `check-go-version` in Makefile to run on WSL. ([443](https://github.com/cucumber/godog/pull/443) - [mxygem](https://github.com/mxygem))

## [v0.12.2]
### Fixed
- Error in `go mod tidy` with `GO111MODULE=off` ([436](https://github.com/cucumber/godog/pull/436) - [vearutop](https://github.com/vearutop))

## [v0.12.1]
### Fixed
- Unintended change of behavior in before step hook ([424](https://github.com/cucumber/godog/pull/424) - [nhatthm](https://github.com/nhatthm))

## [v0.12.0]
### Added
- Support for step definitions without return ([364](https://github.com/cucumber/godog/pull/364) - [titouanfreville](https://github.com/titouanfreville))
- Contextualized hooks for scenarios and steps ([409](https://github.com/cucumber/godog/pull/409) - [vearutop](https://github.com/vearutop))
- Step result status in After hook ([409](https://github.com/cucumber/godog/pull/409) - [vearutop](https://github.com/vearutop))
- Support auto converting doc strings to plain strings ([380](https://github.com/cucumber/godog/pull/380) - [chirino](https://github.com/chirino))
- Use multiple formatters in the same test run ([392](https://github.com/cucumber/godog/pull/392) - [vearutop](https://github.com/vearutop))
- Added `RetrieveFeatures()` method to `godog.TestSuite` ([276](https://github.com/cucumber/godog/pull/276) - [radtriste](https://github.com/radtriste))
- Added support to create custom formatters ([372](https://github.com/cucumber/godog/pull/372) - [leviable](https://github.com/leviable))

### Changed
- Upgraded gherkin-go to v19 and messages-go to v16 ([402](https://github.com/cucumber/godog/pull/402) - [mbow](https://github.com/mbow))
- Generate simpler snippets that use *godog.DocString and *godog.Table ([379](https://github.com/cucumber/godog/pull/379) - [chirino](https://github.com/chirino))

### Deprecated
- `ScenarioContext.BeforeScenario`, use `ScenarioContext.Before` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop](https://github.com/vearutop))
- `ScenarioContext.AfterScenario`, use `ScenarioContext.After` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop](https://github.com/vearutop))
- `ScenarioContext.BeforeStep`, use `ScenarioContext.StepContext().Before` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop](https://github.com/vearutop))
- `ScenarioContext.AfterStep`, use `ScenarioContext.StepContext().After` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop](https://github.com/vearutop))

### Fixed
- Incorrect step definition output for Data Tables ([411](https://github.com/cucumber/godog/pull/411) - [karfrank](https://github.com/karfrank))
- `ScenarioContext.AfterStep` not invoked after a failed case ([409](https://github.com/cucumber/godog/pull/409) - [vearutop](https://github.com/vearutop)))
- Can't execute multiple specific scenarios in the same feature file ([414](https://github.com/cucumber/godog/pull/414) - [vearutop](https://github.com/vearutop)))

## [v0.11.0]
### Added
- Created a simple example for a custom formatter ([330](https://github.com/cucumber/godog/pull/330) - [lonnblad](https://github.com/lonnblad))
- --format junit:result.xml will now write to result.xml ([331](https://github.com/cucumber/godog/pull/331) - [lonnblad](https://github.com/lonnblad))
- Added make commands to create artifacts and upload them to a github release ([333](https://github.com/cucumber/godog/pull/333) - [lonnblad](https://github.com/lonnblad))
- Created release notes and changelog for v0.11.0 ([355](https://github.com/cucumber/godog/pull/355) - [lonnblad](https://github.com/lonnblad))
- Created v0.11.0-rc2 ([362](https://github.com/cucumber/godog/pull/362) - [lonnblad](https://github.com/lonnblad))

### Changed
- Added Cobra for the Command Line Interface ([321](https://github.com/cucumber/godog/pull/321) - [lonnblad](https://github.com/lonnblad))
- Added internal packages for formatters, storage and models ([323](https://github.com/cucumber/godog/pull/323) - [lonnblad](https://github.com/lonnblad))
- Added an internal package for tags filtering ([326](https://github.com/cucumber/godog/pull/326) - [lonnblad](https://github.com/lonnblad))
- Added an internal pkg for the builder ([327](https://github.com/cucumber/godog/pull/327) - [lonnblad](https://github.com/lonnblad))
- Moved the parser code to a new internal pkg ([329](https://github.com/cucumber/godog/pull/329) - [lonnblad](https://github.com/lonnblad))
- Moved StepDefinition to the formatters pkg ([332](https://github.com/cucumber/godog/pull/332) - [lonnblad](https://github.com/lonnblad))
- Removed go1.12 and added go1.15 to CI config ([356](https://github.com/cucumber/godog/pull/356) - [lonnblad](https://github.com/lonnblad))

### Fixed
- Improved the help text of the formatter flag in the run command ([347](https://github.com/cucumber/godog/pull/347) - [lonnblad](https://github.com/lonnblad))
- Removed $GOPATH from the README.md and updated the example ([349](https://github.com/cucumber/godog/pull/349) - [lonnblad](https://github.com/lonnblad))
- Fixed the undefined step definitions help ([350](https://github.com/cucumber/godog/pull/350) - [lonnblad](https://github.com/lonnblad))
- Added a comment regarding running the examples within the $GOPATH ([352](https://github.com/cucumber/godog/pull/352) - [lonnblad](https://github.com/lonnblad))
- doc(FAQ/TestMain): `testing.M.Run()` is optional ([353](https://github.com/cucumber/godog/pull/353) - [hansbogert](https://github.com/hansbogert))
- Made a fix for the unstable Randomize Run tests ([354](https://github.com/cucumber/godog/pull/354) - [lonnblad](https://github.com/lonnblad))
- Fixed an issue when go test is parsing command-line flags ([359](https://github.com/cucumber/godog/pull/359) - [lonnblad](https://github.com/lonnblad))
- Make pickleStepIDs unique accross multiple paths ([366](https://github.com/cucumber/godog/pull/366) - [rickardenglund](https://github.com/rickardenglund))

### Removed
- Removed deprecated code ([322](https://github.com/cucumber/godog/pull/322) - [lonnblad](https://github.com/lonnblad))

## [v0.10.0]
### Added
- Added concurrency support to the pretty formatter ([275](https://github.com/cucumber/godog/pull/275) - [lonnblad](https://github.com/lonnblad))
- Added concurrency support to the events formatter ([274](https://github.com/cucumber/godog/pull/274) - [lonnblad](https://github.com/lonnblad))
- Added concurrency support to the cucumber formatter ([273](https://github.com/cucumber/godog/pull/273) - [lonnblad](https://github.com/lonnblad))
- Added an example for how to use assertion pkgs like testify with godog ([289](https://github.com/cucumber/godog/pull/289) - [lonnblad](https://github.com/lonnblad))
- Added the new TestSuiteInitializer and ScenarioInitializer ([294](https://github.com/cucumber/godog/pull/294) - [lonnblad](https://github.com/lonnblad))
- Added an in-mem storage for pickles ([304](https://github.com/cucumber/godog/pull/304) - [lonnblad](https://github.com/lonnblad))
- Added Pickle and PickleStep results to the in-mem storage ([305](https://github.com/cucumber/godog/pull/305) - [lonnblad](https://github.com/lonnblad))
- Added features to the in-mem storage ([306](https://github.com/cucumber/godog/pull/306) - [lonnblad](https://github.com/lonnblad))
- Broke out some code from massive files into new files ([307](https://github.com/cucumber/godog/pull/307) - [lonnblad](https://github.com/lonnblad))
- Added support for concurrent scenarios ([311](https://github.com/cucumber/godog/pull/311) - [lonnblad](https://github.com/lonnblad))

### Changed
- Broke out snippets gen and added sorting on method name ([271](https://github.com/cucumber/godog/pull/271) - [lonnblad](https://github.com/lonnblad))
- Updated so that we run all tests concurrent now ([278](https://github.com/cucumber/godog/pull/278) - [lonnblad](https://github.com/lonnblad))
- Moved fmt tests to a godog_test pkg and restructured the fmt output tests ([295](https://github.com/cucumber/godog/pull/295) - [lonnblad](https://github.com/lonnblad))
- Moved builder tests to a godog_test pkg ([296](https://github.com/cucumber/godog/pull/296) - [lonnblad](https://github.com/lonnblad))
- Made the builder tests run in parallel ([298](https://github.com/cucumber/godog/pull/298) - [lonnblad](https://github.com/lonnblad))
- Refactored suite_context.go ([300](https://github.com/cucumber/godog/pull/300) - [lonnblad](https://github.com/lonnblad))
- Added better testing of the Context Initializers and TestSuite{}.Run() ([301](https://github.com/cucumber/godog/pull/301) - [lonnblad](https://github.com/lonnblad))
- Updated the README.md ([302](https://github.com/cucumber/godog/pull/302) - [lonnblad](https://github.com/lonnblad))
- Unexported some exported properties in unexported structs ([303](https://github.com/cucumber/godog/pull/303) - [lonnblad](https://github.com/lonnblad))
- Refactored some states in the formatters and feature struct ([310](https://github.com/cucumber/godog/pull/310) - [lonnblad](https://github.com/lonnblad))

### Deprecated
- Deprecated SuiteContext and ConcurrentFormatter ([314](https://github.com/cucumber/godog/pull/314) - [lonnblad](https://github.com/lonnblad))

### Fixed
- Fixed failing builder tests due to the v0.9.0 change ([lonnblad](https://github.com/lonnblad))
- Update paths to screenshots for examples ([270](https://github.com/cucumber/godog/pull/270) - [leviable](https://github.com/leviable))
- Made progress formatter verification a bit more accurate ([lonnblad](https://github.com/lonnblad))
- Added comparison between single and multi threaded runs ([272](https://github.com/cucumber/godog/pull/272) - [lonnblad](https://github.com/lonnblad))
- Fixed issue with empty feature file causing nil pointer deref ([288](https://github.com/cucumber/godog/pull/288) - [lonnblad](https://github.com/lonnblad))
- Updated linting checks in circleci config and fixed linting issues ([290](https://github.com/cucumber/godog/pull/290) - [lonnblad](https://github.com/lonnblad))
- Readded some legacy doc for FeatureContext ([297](https://github.com/cucumber/godog/pull/297) - [lonnblad](https://github.com/lonnblad))
- Fixed an issue with calculating time for junit testsuite ([308](https://github.com/cucumber/godog/pull/308) - [lonnblad](https://github.com/lonnblad))
- Fixed so that we don't execute features with zero scenarios ([315](https://github.com/cucumber/godog/pull/315) - [lonnblad](https://github.com/lonnblad))
- Fixed the broken --random flag ([317](https://github.com/cucumber/godog/pull/317) - [lonnblad](https://github.com/lonnblad))

### Removed
- Removed pre go112 build code ([293](https://github.com/cucumber/godog/pull/293) - [lonnblad](https://github.com/lonnblad))
- Removed the deprecated feature hooks ([312](https://github.com/cucumber/godog/pull/312) - [lonnblad](https://github.com/lonnblad))

## [0.9.0]
### Changed
- Run godog features in CircleCI in strict mode ([mxygem](https://github.com/mxygem))
- Removed TestMain call in `suite_test.go` for CI. ([mxygem](https://github.com/mxygem))
- Migrated to [gherkin-go - v11.0.0](https://github.com/cucumber/gherkin-go/releases/tag/v11.0.0). ([240](https://github.com/cucumber/godog/pull/240) - [lonnblad](https://github.com/lonnblad))

### Fixed
- Fixed the time attributes in the JUnit formatter. ([232](https://github.com/cucumber/godog/pull/232) - [lonnblad](https://github.com/lonnblad))
- Re enable custom formatters. ([238](https://github.com/cucumber/godog/pull/238) - [ericmcbride](https://github.com/ericmcbride))
- Added back suite_test.go ([mxygem](https://github.com/mxygem))
- Normalise module paths for use on Windows ([242](https://github.com/cucumber/godog/pull/242) - [gjtaylor](https://github.com/gjtaylor))
- Fixed panic in indenting function `s` ([247](https://github.com/cucumber/godog/pull/247) - [titouanfreville](https://github.com/titouanfreville))
- Fixed wrong version in API example ([263](https://github.com/cucumber/godog/pull/263) - [denis-trofimov](https://github.com/denis-trofimov))

## [0.8.1]
### Added
- Link in Readme to the Slack community. ([210](https://github.com/cucumber/godog/pull/210) - [smikulcik](https://github.com/smikulcik))
- Added run tests for Cucumber formatting. ([214](https://github.com/cucumber/godog/pull/214), [216](https://github.com/cucumber/godog/pull/216) - [lonnblad](https://github.com/lonnblad))

### Changed
- Renamed the `examples` directory to `_examples`, removing dependencies from the Go module ([218](https://github.com/cucumber/godog/pull/218) - [axw](https://github.com/axw))

### Fixed
- Find/Replaced references to DATA-DOG/godog -> cucumber/godog for docs. ([209](https://github.com/cucumber/godog/pull/209) - [smikulcik](https://github.com/smikulcik))
- Fixed missing links in changelog to be correctly included! ([mxygem](https://github.com/mxygem))

## [0.8.0]
### Added
- Added initial CircleCI config. ([mxygem](https://github.com/mxygem))
- Added concurrency support for JUnit formatting ([lonnblad](https://github.com/lonnblad))

### Changed
- Changed code references to DATA-DOG/godog to cucumber/godog to help get things building correctly. ([mxygem](https://github.com/mxygem))

[v0.15.0]: https://github.com/cucumber/godog/compare/v0.14.1...v0.15.0
[v0.14.1]: https://github.com/cucumber/godog/compare/v0.14.0...v0.14.1
[v0.14.0]: https://github.com/cucumber/godog/compare/v0.13.0...v0.14.0
[v0.13.0]: https://github.com/cucumber/godog/compare/v0.12.6...v0.13.0
[v0.12.6]: https://github.com/cucumber/godog/compare/v0.12.5...v0.12.6
[v0.12.5]: https://github.com/cucumber/godog/compare/v0.12.4...v0.12.5
[v0.12.4]: https://github.com/cucumber/godog/compare/v0.12.3...v0.12.4
[v0.12.3]: https://github.com/cucumber/godog/compare/v0.12.2...v0.12.3
[v0.12.2]: https://github.com/cucumber/godog/compare/v0.12.1...v0.12.2
[v0.12.1]: https://github.com/cucumber/godog/compare/v0.12.0...v0.12.1
[v0.12.0]: https://github.com/cucumber/godog/compare/v0.11.0...v0.12.0
[v0.11.0]: https://github.com/cucumber/godog/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/cucumber/godog/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/cucumber/godog/compare/v0.8.1...v0.9.0
[0.8.1]: https://github.com/cucumber/godog/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/cucumber/godog/compare/v0.7.13...v0.8.0
