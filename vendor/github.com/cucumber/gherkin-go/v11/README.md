[![Build Status](https://secure.travis-ci.org/cucumber/gherkin-go.svg)](http://travis-ci.org/cucumber/gherkin-go) [![GoDoc](https://godoc.org/github.com/cucumber/gherkin-go?status.svg)](http://godoc.org/github.com/cucumber/gherkin-go) [![Go Report Card](https://goreportcard.com/badge/github.com/cucumber/gherkin-go)](https://goreportcard.com/report/github.com/cucumber/gherkin-go)

Gherkin parser/compiler for Go. Please see [Gherkin](https://github.com/cucumber/gherkin) for details.

## Building

You need Go installed (obviously). You also need to make sure your `PATH`
points to where Go installs packages:

```bash
# Add go bin to path
export PATH=$(go env GOPATH)/bin:${PATH}
```

Now build it:

```
make .dist
```

You should have cross-compiled binaries in `./dist/`.

## Compress binaries

You need [upx](https://upx.github.io/) installed.

```
make .dist
make .dist-compressed
```

Your `./dist_compressed/` directory should now have compressed binaries.
Compression fails for some binaries, so you likely won't have a full set.

The build copies the successfully compressed binaries back to `./dist/`.

