# CHANGE LOG

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org).

This document is formatted according to the principles of [Keep A CHANGELOG](http://keepachangelog.com).

---

## [v0.12.2]

### Fixed

- Error in `go mod tidy` with `GO111MODULE=off` ([436]((https://github.com/cucumber/godog/pull/436) - [vearutop])

## [v0.12.1]

### Fixed

- Unintended change of behavior in before step hook ([424](https://github.com/cucumber/godog/pull/424) - [nhatthm])

## [v0.12.0]

### Added

- Support for step definitions without return ([364](https://github.com/cucumber/godog/pull/364) - [titouanfreville])
- Contextualized hooks for scenarios and steps ([409](https://github.com/cucumber/godog/pull/409) - [vearutop])
- Step result status in After hook ([409](https://github.com/cucumber/godog/pull/409) - [vearutop])
- Support auto converting doc strings to plain strings ([380](https://github.com/cucumber/godog/pull/380) - [chirino])
- Use multiple formatters in the same test run ([392](https://github.com/cucumber/godog/pull/392) - [vearutop])
- Added `RetrieveFeatures()` method to `godog.TestSuite` ([276](https://github.com/cucumber/godog/pull/276) - [radtriste])
- Added support to create custom formatters ([372](https://github.com/cucumber/godog/pull/372) - [leviable])

### Changed

- Upgraded gherkin-go to v19 and messages-go to v16 ([402](https://github.com/cucumber/godog/pull/402) - [mbow])
- Generate simpler snippets that use *godog.DocString and *godog.Table ([379](https://github.com/cucumber/godog/pull/379)) - [chirino])

### Deprecated

- `ScenarioContext.BeforeScenario`, use `ScenarioContext.Before` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop])
- `ScenarioContext.AfterScenario`, use `ScenarioContext.After` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop])
- `ScenarioContext.BeforeStep`, use `ScenarioContext.StepContext().Before` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop])
- `ScenarioContext.AfterStep`, use `ScenarioContext.StepContext().After` ([409](https://github.com/cucumber/godog/pull/409)) - [vearutop])

### Removed

### Fixed

- Incorrect step definition output for Data Tables ([411](https://github.com/cucumber/godog/pull/411) - [karfrank])
- `ScenarioContext.AfterStep` not invoked after a failed case (([409](https://github.com/cucumber/godog/pull/409)) - [vearutop]))
- Can't execute multiple specific scenarios in the same feature file (([414](https://github.com/cucumber/godog/pull/414)) - [vearutop]))

## [v0.11.0]

### Added

- Created a simple example for a custom formatter ([330](https://github.com/cucumber/godog/pull/330) - [lonnblad])
- --format junit:result.xml will now write to result.xml ([331](https://github.com/cucumber/godog/pull/331) - [lonnblad])
- Added make commands to create artifacts and upload them to a github release ([333](https://github.com/cucumber/godog/pull/333) - [lonnblad])
- Created release notes and changelog for v0.11.0 ([355](https://github.com/cucumber/godog/pull/355) - [lonnblad])
- Created v0.11.0-rc2 ([362](https://github.com/cucumber/godog/pull/362) - [lonnblad])

### Changed

- Added Cobra for the Command Line Interface ([321](https://github.com/cucumber/godog/pull/321) - [lonnblad])
- Added internal packages for formatters, storage and models ([323](https://github.com/cucumber/godog/pull/323) - [lonnblad])
- Added an internal package for tags filtering ([326](https://github.com/cucumber/godog/pull/326) - [lonnblad])
- Added an internal pkg for the builder ([327](https://github.com/cucumber/godog/pull/327) - [lonnblad])
- Moved the parser code to a new internal pkg ([329](https://github.com/cucumber/godog/pull/329) - [lonnblad])
- Moved StepDefinition to the formatters pkg ([332](https://github.com/cucumber/godog/pull/332) - [lonnblad])
- Removed go1.12 and added go1.15 to CI config ([356](https://github.com/cucumber/godog/pull/356) - [lonnblad])

### Deprecated

### Removed

- Removed deprecated code ([322](https://github.com/cucumber/godog/pull/322) - [lonnblad])

### Fixed

- Improved the help text of the formatter flag in the run command ([347](https://github.com/cucumber/godog/pull/347) - [lonnblad])
- Removed $GOPATH from the README.md and updated the example ([349](https://github.com/cucumber/godog/pull/349) - [lonnblad])
- Fixed the undefined step definitions help ([350](https://github.com/cucumber/godog/pull/350) - [lonnblad])
- Added a comment regarding running the examples within the $GOPATH ([352](https://github.com/cucumber/godog/pull/352) - [lonnblad])
- doc(FAQ/TestMain): `testing.M.Run()` is optional ([353](https://github.com/cucumber/godog/pull/353) - [hansbogert])
- Made a fix for the unstable Randomize Run tests ([354](https://github.com/cucumber/godog/pull/354) - [lonnblad])
- Fixed an issue when go test is parsing command-line flags ([359](https://github.com/cucumber/godog/pull/359) - [lonnblad])
- Make pickleStepIDs unique accross multiple paths ([366](https://github.com/cucumber/godog/pull/366) - [rickardenglund])


## [v0.10.0]

### Added

- Added concurrency support to the pretty formatter ([275](https://github.com/cucumber/godog/pull/275) - [lonnblad])
- Added concurrency support to the events formatter ([274](https://github.com/cucumber/godog/pull/274) - [lonnblad])
- Added concurrency support to the cucumber formatter ([273](https://github.com/cucumber/godog/pull/273) - [lonnblad])
- Added an example for how to use assertion pkgs like testify with godog ([289](https://github.com/cucumber/godog/pull/289) - [lonnblad])
- Added the new TestSuiteInitializer and ScenarioInitializer ([294](https://github.com/cucumber/godog/pull/294) - [lonnblad])
- Added an in-mem storage for pickles ([304](https://github.com/cucumber/godog/pull/304) - [lonnblad])
- Added Pickle and PickleStep results to the in-mem storage ([305](https://github.com/cucumber/godog/pull/305) - [lonnblad])
- Added features to the in-mem storage ([306](https://github.com/cucumber/godog/pull/306) - [lonnblad])
- Broke out some code from massive files into new files ([307](https://github.com/cucumber/godog/pull/307) - [lonnblad])
- Added support for concurrent scenarios ([311](https://github.com/cucumber/godog/pull/311) - [lonnblad])

### Changed

- Broke out snippets gen and added sorting on method name ([271](https://github.com/cucumber/godog/pull/271) - [lonnblad])
- Updated so that we run all tests concurrent now ([278](https://github.com/cucumber/godog/pull/278) - [lonnblad])
- Moved fmt tests to a godog_test pkg and restructured the fmt output tests ([295](https://github.com/cucumber/godog/pull/295) - [lonnblad])
- Moved builder tests to a godog_test pkg ([296](https://github.com/cucumber/godog/pull/296) - [lonnblad])
- Made the builder tests run in parallel ([298](https://github.com/cucumber/godog/pull/298) - [lonnblad])
- Refactored suite_context.go ([300](https://github.com/cucumber/godog/pull/300) - [lonnblad])
- Added better testing of the Context Initializers and TestSuite{}.Run() ([301](https://github.com/cucumber/godog/pull/301) - [lonnblad])
- Updated the README.md ([302](https://github.com/cucumber/godog/pull/302) - [lonnblad])
- Unexported some exported properties in unexported structs ([303](https://github.com/cucumber/godog/pull/303) - [lonnblad])
- Refactored some states in the formatters and feature struct ([310](https://github.com/cucumber/godog/pull/310) - [lonnblad])

### Deprecated

- Deprecated SuiteContext and ConcurrentFormatter ([314](https://github.com/cucumber/godog/pull/314) - [lonnblad])

### Removed

- Removed pre go112 build code ([293](https://github.com/cucumber/godog/pull/293) - [lonnblad])
- Removed the deprecated feature hooks ([312](https://github.com/cucumber/godog/pull/312) - [lonnblad])

### Fixed

- Fixed failing builder tests due to the v0.9.0 change ([lonnblad])
- Update paths to screenshots for examples ([270](https://github.com/cucumber/godog/pull/270) - [leviable])
- Made progress formatter verification a bit more accurate ([lonnblad])
- Added comparison between single and multi threaded runs ([272](https://github.com/cucumber/godog/pull/272) - [lonnblad])
- Fixed issue with empty feature file causing nil pointer deref ([288](https://github.com/cucumber/godog/pull/288) - [lonnblad])
- Updated linting checks in circleci config and fixed linting issues ([290](https://github.com/cucumber/godog/pull/290) - [lonnblad])
- Readded some legacy doc for FeatureContext ([297](https://github.com/cucumber/godog/pull/297) - [lonnblad])
- Fixed an issue with calculating time for junit testsuite ([308](https://github.com/cucumber/godog/pull/308) - [lonnblad])
- Fixed so that we don't execute features with zero scenarios ([315](https://github.com/cucumber/godog/pull/315) - [lonnblad])
- Fixed the broken --random flag ([317](https://github.com/cucumber/godog/pull/317) - [lonnblad])

## [0.9.0]

### Added

### Changed

- Run godog features in CircleCI in strict mode ([jaysonesmith])
- Removed TestMain call in `suite_test.go` for CI. ([jaysonesmith])
- Migrated to [gherkin-go - v11.0.0](https://github.com/cucumber/gherkin-go/releases/tag/v11.0.0). ([240](https://github.com/cucumber/godog/pull/240) - [lonnblad])

### Deprecated

### Removed

### Fixed

- Fixed the time attributes in the JUnit formatter. ([232](https://github.com/cucumber/godog/pull/232) - [lonnblad])
- Re enable custom formatters. ([238](https://github.com/cucumber/godog/pull/238) - [ericmcbride])
- Added back suite_test.go ([jaysonesmith])
- Normalise module paths for use on Windows ([242](https://github.com/cucumber/godog/pull/242) - [gjtaylor])
- Fixed panic in indenting function `s` ([247](https://github.com/cucumber/godog/pull/247) - [titouanfreville])
- Fixed wrong version in API example ([263](https://github.com/cucumber/godog/pull/263) - [denis-trofimov])

## [0.8.1]

### Added

- Link in Readme to the Slack community. ([210](https://github.com/cucumber/godog/pull/210) - [smikulcik])
- Added run tests for Cucumber formatting. ([214](https://github.com/cucumber/godog/pull/214), [216](https://github.com/cucumber/godog/pull/216) - [lonnblad])

### Changed

- Renamed the `examples` directory to `_examples`, removing dependencies from the Go module ([218](https://github.com/cucumber/godog/pull/218) - [axw])

### Deprecated

### Removed

### Fixed

- Find/Replaced references to DATA-DOG/godog -> cucumber/godog for docs. ([209](https://github.com/cucumber/godog/pull/209) - [smikulcik])
- Fixed missing links in changelog to be correctly included! ([jaysonesmith])

## [0.8.0]

### Added

- Added initial CircleCI config. ([jaysonesmith])
- Added concurrency support for JUnit formatting ([lonnblad])

### Changed

- Changed code references to DATA-DOG/godog to cucumber/godog to help get things building correctly. ([jaysonesmith])

### Deprecated

### Removed

### Fixed

<!-- Releases -->

[v0.12.2]: https://github.com/cucumber/godog/compare/v0.12.1...v0.12.2
[v0.12.1]: https://github.com/cucumber/godog/compare/v0.12.0...v0.12.1
[v0.12.0]: https://github.com/cucumber/godog/compare/v0.11.0...v0.12.0
[v0.11.0]: https://github.com/cucumber/godog/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/cucumber/godog/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/cucumber/godog/compare/v0.8.1...v0.9.0
[0.8.1]: https://github.com/cucumber/godog/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/cucumber/godog/compare/v0.7.13...v0.8.0

<!-- Contributors -->

[axw]: https://github.com/axw
[jaysonesmith]: https://github.com/jaysonesmith
[lonnblad]: https://github.com/lonnblad
[smikulcik]: https://github.com/smikulcik
[ericmcbride]: https://github.com/ericmcbride
[gjtaylor]: https://github.com/gjtaylor
[titouanfreville]: https://github.com/titouanfreville
[denis-trofimov]: https://github.com/denis-trofimov
[leviable]: https://github.com/leviable
[hansbogert]: https://github.com/hansbogert
[rickardenglund]: https://github.com/rickardenglund
[mbow]: https://github.com/mbow
[vearutop]: https://github.com/vearutop
[chirino]: https://github.com/chirino
[radtriste]: https://github.com/radtriste
[karfrank]: https://github.com/karfrank
[nhatthm]: https://github.com/nhatthm
