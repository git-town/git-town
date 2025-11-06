# displaywidth

A high-performance Go package for measuring the monospace display width of strings, UTF-8 bytes, and runes.

[![Documentation](https://pkg.go.dev/badge/github.com/clipperhouse/displaywidth.svg)](https://pkg.go.dev/github.com/clipperhouse/displaywidth)
[![Test](https://github.com/clipperhouse/displaywidth/actions/workflows/gotest.yml/badge.svg)](https://github.com/clipperhouse/displaywidth/actions/workflows/gotest.yml)
[![Fuzz](https://github.com/clipperhouse/displaywidth/actions/workflows/gofuzz.yml/badge.svg)](https://github.com/clipperhouse/displaywidth/actions/workflows/gofuzz.yml)
## Install
```bash
go get github.com/clipperhouse/displaywidth
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/clipperhouse/displaywidth"
)

func main() {
    width := displaywidth.String("Hello, 世界!")
    fmt.Println(width)

    width = displaywidth.Bytes([]byte("🌍"))
    fmt.Println(width)

    width = displaywidth.Rune('🌍')
    fmt.Println(width)
}
```

For most purposes, you should use the `String` or `Bytes` methods.

### Options

You can specify East Asian Width and Strict Emoji Neutral settings. If
unspecified, the default is `EastAsianWidth: false, StrictEmojiNeutral: true`.

```go
options := displaywidth.Options{
    EastAsianWidth:     true,
    StrictEmojiNeutral: false,
}

width := options.String("Hello, 世界!")
fmt.Println(width)
```

## Details

This package implements the Unicode East Asian Width standard
([UAX #11](https://www.unicode.org/reports/tr11/)), and handles
[version selectors](https://en.wikipedia.org/wiki/Variation_Selectors_(Unicode_block)),
and [regional indicator pairs](https://en.wikipedia.org/wiki/Regional_indicator_symbol)
(flags). It operates on bytes without decoding runes for better performance.

## Prior Art

[mattn/go-runewidth](https://github.com/mattn/go-runewidth)

[rivo/uniseg](https://github.com/rivo/uniseg)

[x/text/width](https://pkg.go.dev/golang.org/x/text/width)

[x/text/internal/triegen](https://pkg.go.dev/golang.org/x/text/internal/triegen)

## Benchmarks

```bash
cd comparison
go test -bench=. -benchmem
```

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/displaywidth/comparison
cpu: Apple M2
BenchmarkStringDefault/clipperhouse/displaywidth-8           11124 ns/op	 151.66 MB/s	  0 B/op	   0 allocs/op
BenchmarkStringDefault/mattn/go-runewidth-8                  14209 ns/op	 118.73 MB/s	  0 B/op	   0 allocs/op
BenchmarkStringDefault/rivo/uniseg-8                         19287 ns/op	  87.47 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_EAW/clipperhouse/displaywidth-8              11134 ns/op	 151.52 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_EAW/mattn/go-runewidth-8                     23793 ns/op	  70.90 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_EAW/rivo/uniseg-8                            19593 ns/op	  86.10 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_StrictEmoji/clipperhouse/displaywidth-8      11124 ns/op	 151.65 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_StrictEmoji/mattn/go-runewidth-8             14310 ns/op	 117.89 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_StrictEmoji/rivo/uniseg-8                    19699 ns/op	  85.64 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_ASCII/clipperhouse/displaywidth-8             1107 ns/op	 115.64 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                    1174 ns/op	 109.04 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                           1582 ns/op	  80.93 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_Unicode/clipperhouse/displaywidth-8            959.7 ns/op	 138.59 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_Unicode/mattn/go-runewidth-8                  1422 ns/op	  93.53 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_Unicode/rivo/uniseg-8                         2032 ns/op	  65.47 MB/s	  0 B/op	   0 allocs/op
BenchmarkStringWidth_Emoji/clipperhouse/displaywidth-8        3230 ns/op	 224.17 MB/s	  0 B/op	   0 allocs/op
BenchmarkStringWidth_Emoji/mattn/go-runewidth-8               4826 ns/op	 150.03 MB/s	  0 B/op	   0 allocs/op
BenchmarkStringWidth_Emoji/rivo/uniseg-8                      6582 ns/op	 109.99 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_Mixed/clipperhouse/displaywidth-8             4094 ns/op	 123.83 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                    4612 ns/op	 109.92 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                           6312 ns/op	  80.32 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_ControlChars/clipperhouse/displaywidth-8       346.7 ns/op	  95.19 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_ControlChars/mattn/go-runewidth-8              365.0 ns/op	  90.42 MB/s	  0 B/op	   0 allocs/op
BenchmarkString_ControlChars/rivo/uniseg-8                     408.9 ns/op	  80.70 MB/s	  0 B/op	   0 allocs/op
```

I use a similar technique in [this grapheme cluster library](https://github.com/clipperhouse/uax29).

## Compatibility

`clipperhouse/displaywidth`, `mattn/go-runewidth`, and `rivo/uniseg` should give the
same outputs for real-world text. See [comparison/README.md](comparison/README.md).
