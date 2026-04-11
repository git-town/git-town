package displaywidth

import (
	"unicode/utf8"

	"github.com/clipperhouse/stringish"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

// String calculates the display width of a string
// using the [DefaultOptions]
func String(s string) int {
	return DefaultOptions.String(s)
}

// Bytes calculates the display width of a []byte
// using the [DefaultOptions]
func Bytes(s []byte) int {
	return DefaultOptions.Bytes(s)
}

func Rune(r rune) int {
	return DefaultOptions.Rune(r)
}

type Options struct {
	EastAsianWidth     bool
	StrictEmojiNeutral bool
}

var DefaultOptions = Options{
	EastAsianWidth:     false,
	StrictEmojiNeutral: true,
}

// String calculates the display width of a string
// for the given options
func (options Options) String(s string) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromString(s)
	for g.Next() {
		// The first character in the grapheme cluster determines the width;
		// we use lookupProperties which can consider immediate VS15/VS16.
		props := lookupProperties(g.Value())
		total += props.width(options)
	}
	return total
}

// BytesOptions calculates the display width of a []byte
// for the given options
func (options Options) Bytes(s []byte) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromBytes(s)
	for g.Next() {
		// The first character in the grapheme cluster determines the width;
		// we use lookupProperties which can consider immediate VS15/VS16.
		props := lookupProperties(g.Value())
		total += props.width(options)
	}
	return total
}

func (options Options) Rune(r rune) int {
	// Fast path for ASCII
	if r < utf8.RuneSelf {
		if isASCIIControl(byte(r)) {
			// Control (0x00-0x1F) and DEL (0x7F)
			return 0
		}
		// ASCII printable (0x20-0x7E)
		return 1
	}

	// Surrogates (U+D800-U+DFFF) are invalid UTF-8 and have zero width
	// Other packages might turn them into the replacement character (U+FFFD)
	// in which case, we won't see it.
	if r >= 0xD800 && r <= 0xDFFF {
		return 0
	}

	// Stack-allocated to avoid heap allocation
	var buf [4]byte // UTF-8 is at most 4 bytes
	n := utf8.EncodeRune(buf[:], r)
	// Skip the grapheme iterator and directly lookup properties
	props := lookupProperties(buf[:n])
	return props.width(options)
}

func isASCIIControl(b byte) bool {
	return b < 0x20 || b == 0x7F
}

const defaultWidth = 1

// is returns true if the property flag is set
func (p property) is(flag property) bool {
	return p&flag != 0
}

// lookupProperties returns the properties for the first character in a string
func lookupProperties[T stringish.Interface](s T) property {
	if len(s) == 0 {
		return 0
	}

	b := s[0]
	if isASCIIControl(b) {
		return _ZeroWidth
	}

	l := len(s)

	if b < utf8.RuneSelf { // Single-byte ASCII
		// Check for variation selector after ASCII (e.g., keycap sequences like 1️⃣)
		var p property
		if l >= 4 {
			// Create a subslice to help the compiler eliminate bounds checks
			vs := s[1:4]
			if vs[0] == 0xEF && vs[1] == 0xB8 {
				switch vs[2] {
				case 0x8E:
					p |= _VS15
				case 0x8F:
					p |= _VS16
				}
			}
		}
		return p // ASCII characters are width 1 by default, or 2 with VS16
	}

	// Regional indicator pair (flag) - detect early before trie lookup.
	// Formed by two Regional Indicator symbols (U+1F1E6–U+1F1FF),
	// each encoded as F0 9F 87 A6–BF. Always width 2, no trie lookup needed.
	if l >= 8 {
		// Create a subslice to help the compiler eliminate bounds checks
		ri := s[:8]
		if ri[0] == 0xF0 &&
			ri[1] == 0x9F &&
			ri[2] == 0x87 {
			b3 := ri[3]
			if b3 >= 0xA6 && b3 <= 0xBF &&
				ri[4] == 0xF0 &&
				ri[5] == 0x9F &&
				ri[6] == 0x87 {
				b7 := ri[7]
				if b7 >= 0xA6 && b7 <= 0xBF {
					return _RI_PAIR
				}
			}
		}
	}

	props, size := lookup(s)
	p := property(props)

	// Variation Selectors
	if size > 0 && l >= size+3 {
		// Create a subslice to help the compiler eliminate bounds checks
		vs := s[size : size+3]
		if vs[0] == 0xEF && vs[1] == 0xB8 {
			switch vs[2] {
			case 0x8E:
				p |= _VS15
			case 0x8F:
				p |= _VS16
			}
		}
	}

	return p
}

// width determines the display width of a character based on its properties
// and configuration options
func (p property) width(options Options) int {
	if p == 0 {
		// Character not in trie, use default behavior
		return defaultWidth
	}

	if p.is(_ZeroWidth) {
		return 0
	}

	// Explicit presentation overrides from VS come first.
	if p.is(_VS16) {
		return 2
	}
	if p.is(_VS15) {
		return 1
	}

	// Regional indicator pair (flag) grapheme cluster
	// returns 1 under StrictEmojiNeutral=false, which
	// is compatible with go-runewidth & uniseg.
	if p.is(_RI_PAIR) && options.StrictEmojiNeutral {
		return 2
	}

	if options.EastAsianWidth {
		if p.is(_East_Asian_Ambiguous) {
			return 2
		}
		if p.is(_East_Asian_Ambiguous|_Emoji) && !options.StrictEmojiNeutral {
			return 2
		}
	}

	if p.is(_East_Asian_Full_Wide) {
		return 2
	}

	// Default width for all other characters
	return defaultWidth
}
