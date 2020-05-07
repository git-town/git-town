# CHANGE LOG

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org).

This document is formatted according to the principles of [Keep A CHANGELOG](http://keepachangelog.com).

----

## [Unreleased]

### Added

### Changed

- Run godog features in CircleCI in strict mode ([jaysonesmith])
- Removed TestMain call in `suite_test.go` for CI. ([jaysonesmith])

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
[Unreleased]: https://github.com/cucumber/cucumber/compare/godog/v0.8.1...master
[0.8.0]:      https://github.com/cucumber/cucumber/compare/godog/v0.8.0...godog/v0.8.1
[0.8.0]:      https://github.com/cucumber/cucumber/compare/godog/v0.7.13...godog/v0.8.0

<!-- Contributors -->
[axw]:              https://github.com/axw
[jaysonesmith]:     https://github.com/jaysonesmith
[lonnblad]:         https://github.com/lonnblad
[smikulcik]:        https://github.com/smikulcik
[ericmcbride]:      https://github.com/ericmcbride
[gjtaylor]:         https://github.com/gjtaylor
[titouanfreville]:  https://github.com/titouanfreville
[denis-trofimov]:   https://github.com/denis-trofimov
