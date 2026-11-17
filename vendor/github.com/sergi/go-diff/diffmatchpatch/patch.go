// Copyright (c) 2012-2016 The go-diff authors. All rights reserved.
// https://github.com/sergi/go-diff
// See the included LICENSE file for license details.
//
// go-diff is a Go implementation of Google's Diff, Match, and Patch library
// Original library is Copyright (c) 2006 Google Inc.
// http://code.google.com/p/google-diff-match-patch/

package diffmatchpatch

import (
	"bytes"
	"errors"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// alignRuneStart backs up byte index idx to the nearest UTF-8 rune start byte.
// idx <= 0 or idx >= len(s) is returned unchanged.
//
// Motivation: Google's reference JS diff-match-patch uses string.substring
// (UTF-16 code unit) for all slicing — surrogate-pair-middle splits are
// well-defined. This Go port slices by UTF-8 byte index (text[a:b]). When
// MatchMain's Bitap returns a byte position that lands inside a multi-byte
// character (CJK 3B / em dash 3B / emoji 4B), the resulting slice is invalid
// UTF-8 — subsequent MatchMain calls miss anchors, and downstream byte slices
// can panic with "slice bounds out of range" (issue #132).
//
// Fix: align every byte slice index back to the nearest valid rune boundary
// before slicing. The operation costs at most 1-3 bytes of "shrinkage" per
// slice point but guarantees the resulting string is always valid UTF-8 and
// no MatchMain / DiffMain pass sees a continuation-byte-prefixed string.
func alignRuneStart(s string, idx int) int {
	if idx <= 0 || idx >= len(s) {
		return idx
	}
	for idx > 0 && !utf8.RuneStart(s[idx]) {
		idx--
	}
	return idx
}

// alignRuneEnd moves byte index idx forward to the next UTF-8 rune start.
// idx <= 0 or idx >= len(s) is returned unchanged.
func alignRuneEnd(s string, idx int) int {
	for idx > 0 && idx < len(s) && !utf8.RuneStart(s[idx]) {
		idx++
	}
	return idx
}

// Patch represents one patch operation.
type Patch struct {
	diffs   []Diff
	Start1  int
	Start2  int
	Length1 int
	Length2 int
}

// String emulates GNU diff's format.
// Header: @@ -382,8 +481,9 @@
// Indices are printed as 1-based, not 0-based.
func (p *Patch) String() string {
	var coords1, coords2 string

	if p.Length1 == 0 {
		coords1 = strconv.Itoa(p.Start1) + ",0"
	} else if p.Length1 == 1 {
		coords1 = strconv.Itoa(p.Start1 + 1)
	} else {
		coords1 = strconv.Itoa(p.Start1+1) + "," + strconv.Itoa(p.Length1)
	}

	if p.Length2 == 0 {
		coords2 = strconv.Itoa(p.Start2) + ",0"
	} else if p.Length2 == 1 {
		coords2 = strconv.Itoa(p.Start2 + 1)
	} else {
		coords2 = strconv.Itoa(p.Start2+1) + "," + strconv.Itoa(p.Length2)
	}

	var text bytes.Buffer
	_, _ = text.WriteString("@@ -" + coords1 + " +" + coords2 + " @@\n")

	// Escape the body of the patch with %xx notation.
	for _, aDiff := range p.diffs {
		switch aDiff.Type {
		case DiffInsert:
			_, _ = text.WriteString("+")
		case DiffDelete:
			_, _ = text.WriteString("-")
		case DiffEqual:
			_, _ = text.WriteString(" ")
		}

		_, _ = text.WriteString(strings.Replace(url.QueryEscape(aDiff.Text), "+", " ", -1))
		_, _ = text.WriteString("\n")
	}

	return unescaper.Replace(text.String())
}

// PatchAddContext increases the context until it is unique, but doesn't let the pattern expand beyond MatchMaxBits.
func (dmp *DiffMatchPatch) PatchAddContext(patch Patch, text string) Patch {
	if len(text) == 0 {
		return patch
	}

	pattern := text[patch.Start2 : patch.Start2+patch.Length1]
	padding := 0

	// Look for the first and last matches of pattern in text.  If two different matches are found, increase the pattern length.
	for strings.Index(text, pattern) != strings.LastIndex(text, pattern) &&
		len(pattern) < dmp.MatchMaxBits-2*dmp.PatchMargin {
		padding += dmp.PatchMargin
		maxStart := max(0, patch.Start2-padding)
		minEnd := min(len(text), patch.Start2+patch.Length1+padding)
		pattern = text[maxStart:minEnd]
	}
	// Add one chunk for good luck.
	padding += dmp.PatchMargin

	// Add the prefix.
	// rune-safety: align prefix start to nearest rune boundary so the slice
	// never produces an invalid-UTF-8 string when the surrounding text contains
	// multi-byte characters.
	prefixStart := alignRuneStart(text, max(0, patch.Start2-padding))
	prefix := text[prefixStart:patch.Start2]
	if len(prefix) != 0 {
		patch.diffs = append([]Diff{Diff{DiffEqual, prefix}}, patch.diffs...)
	}
	// Add the suffix.
	// rune-safety: round suffix end forward to the next rune boundary. Rounding
	// back would shorten the context, and PatchAddPadding would then treat the
	// patch as touching the end of the text.
	suffixEnd := alignRuneEnd(text, min(len(text), patch.Start2+patch.Length1+padding))
	suffix := text[patch.Start2+patch.Length1 : suffixEnd]
	if len(suffix) != 0 {
		patch.diffs = append(patch.diffs, Diff{DiffEqual, suffix})
	}

	// Roll back the start points.
	patch.Start1 -= len(prefix)
	patch.Start2 -= len(prefix)
	// Extend the lengths.
	patch.Length1 += len(prefix) + len(suffix)
	patch.Length2 += len(prefix) + len(suffix)

	return patch
}

// PatchMake computes a list of patches.
func (dmp *DiffMatchPatch) PatchMake(opt ...interface{}) []Patch {
	if len(opt) == 1 {
		diffs, _ := opt[0].([]Diff)
		text1 := dmp.DiffText1(diffs)
		return dmp.PatchMake(text1, diffs)
	} else if len(opt) == 2 {
		text1 := opt[0].(string)
		switch t := opt[1].(type) {
		case string:
			diffs := dmp.DiffMain(text1, t, true)
			if len(diffs) > 2 {
				diffs = dmp.DiffCleanupSemantic(diffs)
				diffs = dmp.DiffCleanupEfficiency(diffs)
			}
			return dmp.PatchMake(text1, diffs)
		case []Diff:
			return dmp.patchMake2(text1, t)
		}
	} else if len(opt) == 3 {
		return dmp.PatchMake(opt[0], opt[2])
	}
	return []Patch{}
}

// patchMake2 computes a list of patches to turn text1 into text2.
// text2 is not provided, diffs are the delta between text1 and text2.
func (dmp *DiffMatchPatch) patchMake2(text1 string, diffs []Diff) []Patch {
	// Check for null inputs not needed since null can't be passed in C#.
	patches := []Patch{}
	if len(diffs) == 0 {
		return patches // Get rid of the null case.
	}

	// DiffMain replaces invalid UTF-8 with the Unicode replacement character.
	// Use that normalized source when it is exactly what the diffs describe,
	// while preserving raw bytes in caller-supplied byte-level diffs.
	if !utf8.ValidString(text1) {
		diffText1 := dmp.DiffText1(diffs)
		if string([]rune(text1)) == diffText1 {
			text1 = diffText1
		}
	}

	patch := Patch{}
	charCount1 := 0 // Number of characters into the text1 string.
	charCount2 := 0 // Number of characters into the text2 string.
	// Start with text1 (prepatchText) and apply the diffs until we arrive at text2 (postpatchText). We recreate the patches one by one to determine context info.
	prepatchText := text1
	postpatchText := text1

	for i, aDiff := range diffs {
		if len(patch.diffs) == 0 && aDiff.Type != DiffEqual {
			// A new patch starts here.
			patch.Start1 = charCount1
			patch.Start2 = charCount2
		}

		switch aDiff.Type {
		case DiffInsert:
			patch.diffs = append(patch.diffs, aDiff)
			patch.Length2 += len(aDiff.Text)
			postpatchText = postpatchText[:charCount2] +
				aDiff.Text + postpatchText[charCount2:]
		case DiffDelete:
			patch.Length1 += len(aDiff.Text)
			patch.diffs = append(patch.diffs, aDiff)
			postpatchText = postpatchText[:charCount2] + postpatchText[charCount2+len(aDiff.Text):]
		case DiffEqual:
			if len(aDiff.Text) <= 2*dmp.PatchMargin &&
				len(patch.diffs) != 0 && i != len(diffs)-1 {
				// Small equality inside a patch.
				patch.diffs = append(patch.diffs, aDiff)
				patch.Length1 += len(aDiff.Text)
				patch.Length2 += len(aDiff.Text)
			}
			if len(aDiff.Text) >= 2*dmp.PatchMargin {
				// Time for a new patch.
				if len(patch.diffs) != 0 {
					patch = dmp.PatchAddContext(patch, prepatchText)
					patches = append(patches, patch)
					patch = Patch{}
					// Unlike Unidiff, our patch lists have a rolling context. http://code.google.com/p/google-diff-match-patch/wiki/Unidiff Update prepatch text & pos to reflect the application of the just completed patch.
					prepatchText = postpatchText
					charCount1 = charCount2
				}
			}
		}

		// Update the current character count.
		if aDiff.Type != DiffInsert {
			charCount1 += len(aDiff.Text)
		}
		if aDiff.Type != DiffDelete {
			charCount2 += len(aDiff.Text)
		}
	}

	// Pick up the leftover patch if not empty.
	if len(patch.diffs) != 0 {
		patch = dmp.PatchAddContext(patch, prepatchText)
		patches = append(patches, patch)
	}

	return patches
}

// PatchDeepCopy returns an array that is identical to a given an array of patches.
func (dmp *DiffMatchPatch) PatchDeepCopy(patches []Patch) []Patch {
	patchesCopy := []Patch{}
	for _, aPatch := range patches {
		patchCopy := Patch{}
		for _, aDiff := range aPatch.diffs {
			patchCopy.diffs = append(patchCopy.diffs, Diff{
				aDiff.Type,
				aDiff.Text,
			})
		}
		patchCopy.Start1 = aPatch.Start1
		patchCopy.Start2 = aPatch.Start2
		patchCopy.Length1 = aPatch.Length1
		patchCopy.Length2 = aPatch.Length2
		patchesCopy = append(patchesCopy, patchCopy)
	}
	return patchesCopy
}

// PatchApply merges a set of patches onto the text.  Returns a patched text, as well as an array of true/false values indicating which patches were applied.
func (dmp *DiffMatchPatch) PatchApply(patches []Patch, text string) (string, []bool) {
	if len(patches) == 0 {
		return text, []bool{}
	}

	// Deep copy the patches so that no changes are made to originals.
	patches = dmp.PatchDeepCopy(patches)

	nullPadding := dmp.PatchAddPadding(patches)
	text = nullPadding + text + nullPadding
	patches = dmp.PatchSplitMax(patches)

	x := 0
	// delta keeps track of the offset between the expected and actual location of the previous patch.  If there are patches expected at positions 10 and 20, but the first patch was found at 12, delta is 2 and the second patch has an effective expected position of 22.
	delta := 0
	results := make([]bool, len(patches))
	for _, aPatch := range patches {
		expectedLoc := aPatch.Start2 + delta
		text1 := dmp.DiffText1(aPatch.diffs)
		var startLoc int
		endLoc := -1
		if len(text1) > dmp.MatchMaxBits {
			// PatchSplitMax will only provide an oversized pattern in the case of a monster delete.
			startLoc = dmp.MatchMain(text, text1[:dmp.MatchMaxBits], expectedLoc)
			if startLoc != -1 {
				endLoc = dmp.MatchMain(text,
					text1[len(text1)-dmp.MatchMaxBits:], expectedLoc+len(text1)-dmp.MatchMaxBits)
				if endLoc == -1 || startLoc >= endLoc {
					// Can't find valid trailing context.  Drop this patch.
					startLoc = -1
				}
			}
		} else {
			startLoc = dmp.MatchMain(text, text1, expectedLoc)
		}
		// rune-safety: Bitap-returned byte position may land in the middle of a
		// multi-byte character. Back it up to the nearest rune start so the
		// subsequent text[startLoc:...] / text[:startLoc] slices never produce
		// invalid UTF-8 (root cause of issue #132 panic).
		if startLoc != -1 {
			startLoc = alignRuneStart(text, startLoc)
		}
		if startLoc == -1 {
			// No match found.  :(
			results[x] = false
			// Subtract the delta for this failed patch from subsequent patches.
			delta -= aPatch.Length2 - aPatch.Length1
		} else {
			// Found a match.  :)
			results[x] = true
			delta = startLoc - expectedLoc
			var text2 string
			var text2End int
			// rune-safety: round the end of text2 forward so it keeps the whole
			// last rune; rounding back would drop it and a deletion of that rune
			// would then silently do nothing.
			if endLoc == -1 {
				text2End = alignRuneEnd(text, int(math.Min(float64(startLoc+len(text1)), float64(len(text)))))
			} else {
				text2End = alignRuneEnd(text, int(math.Min(float64(endLoc+dmp.MatchMaxBits), float64(len(text)))))
			}
			text2 = text[startLoc:text2End]
			if text1 == text2 {
				// Perfect match, just shove the Replacement text in.
				replaceEnd := alignRuneStart(text, startLoc+len(text1))
				text = text[:startLoc] + dmp.DiffText2(aPatch.diffs) + text[replaceEnd:]
			} else {
				// Imperfect match.  Run a diff to get a framework of equivalent indices.
				diffs := dmp.DiffMain(text1, text2, false)
				if len(text1) > dmp.MatchMaxBits && float64(dmp.DiffLevenshtein(diffs))/float64(len(text1)) > dmp.PatchDeleteThreshold {
					// The end points match, but the content is unacceptably bad.
					results[x] = false
				} else {
					diffs = dmp.DiffCleanupSemanticLossless(diffs)
					index1 := 0
					for _, aDiff := range aPatch.diffs {
						if aDiff.Type != DiffEqual {
							// Indices come from text2, but earlier edits in this patch can
							// leave text shorter than that. Clamp them to the end of text,
							// as substring does in the reference implementation.
							index2 := dmp.DiffXIndex(diffs, index1)
							if aDiff.Type == DiffInsert {
								// Insertion
								insertAt := alignRuneStart(text, min(startLoc+index2, len(text)))
								text = text[:insertAt] + aDiff.Text + text[insertAt:]
							} else if aDiff.Type == DiffDelete {
								// Deletion
								startIndex := alignRuneStart(text, min(startLoc+index2, len(text)))
								endIndex := alignRuneStart(text, min(startLoc+dmp.DiffXIndex(diffs, index1+len(aDiff.Text)), len(text)))
								if endIndex < startIndex {
									endIndex = startIndex
								}
								text = text[:startIndex] + text[endIndex:]
							}
						}
						if aDiff.Type != DiffDelete {
							index1 += len(aDiff.Text)
						}
					}
				}
			}
		}
		x++
	}
	// Strip the padding off. Deletions can eat into the padding, so clamp the
	// cut points to the text. An insertion that landed in the padding can
	// straddle a cut point; drop that rune so the result stays valid UTF-8.
	start := alignRuneEnd(text, min(len(nullPadding), len(text)))
	end := max(start, alignRuneStart(text, len(text)-len(nullPadding)))
	text = text[start:end]
	return text, results
}

// PatchAddPadding adds some padding on text start and end so that edges can match something.
// Intended to be called only from within patchApply.
func (dmp *DiffMatchPatch) PatchAddPadding(patches []Patch) string {
	paddingLength := dmp.PatchMargin
	nullPadding := ""
	for x := 1; x <= paddingLength; x++ {
		nullPadding += string(rune(x))
	}

	// Bump all the patches forward.
	for i := range patches {
		patches[i].Start1 += paddingLength
		patches[i].Start2 += paddingLength
	}

	// Add some padding on start of first diff.
	if len(patches[0].diffs) == 0 || patches[0].diffs[0].Type != DiffEqual {
		// Add nullPadding equality.
		patches[0].diffs = append([]Diff{Diff{DiffEqual, nullPadding}}, patches[0].diffs...)
		patches[0].Start1 -= paddingLength // Should be 0.
		patches[0].Start2 -= paddingLength // Should be 0.
		patches[0].Length1 += paddingLength
		patches[0].Length2 += paddingLength
	} else if paddingLength > len(patches[0].diffs[0].Text) {
		// Grow first equality.
		extraLength := paddingLength - len(patches[0].diffs[0].Text)
		patches[0].diffs[0].Text = nullPadding[len(patches[0].diffs[0].Text):] + patches[0].diffs[0].Text
		patches[0].Start1 -= extraLength
		patches[0].Start2 -= extraLength
		patches[0].Length1 += extraLength
		patches[0].Length2 += extraLength
	}

	// Add some padding on end of last diff.
	last := len(patches) - 1
	if len(patches[last].diffs) == 0 || patches[last].diffs[len(patches[last].diffs)-1].Type != DiffEqual {
		// Add nullPadding equality.
		patches[last].diffs = append(patches[last].diffs, Diff{DiffEqual, nullPadding})
		patches[last].Length1 += paddingLength
		patches[last].Length2 += paddingLength
	} else if paddingLength > len(patches[last].diffs[len(patches[last].diffs)-1].Text) {
		// Grow last equality.
		lastDiff := patches[last].diffs[len(patches[last].diffs)-1]
		extraLength := paddingLength - len(lastDiff.Text)
		patches[last].diffs[len(patches[last].diffs)-1].Text += nullPadding[:extraLength]
		patches[last].Length1 += extraLength
		patches[last].Length2 += extraLength
	}

	return nullPadding
}

// PatchSplitMax looks through the patches and breaks up any which are longer than the maximum limit of the match algorithm.
// Intended to be called only from within patchApply.
func (dmp *DiffMatchPatch) PatchSplitMax(patches []Patch) []Patch {
	patchSize := dmp.MatchMaxBits
	for x := 0; x < len(patches); x++ {
		if patches[x].Length1 <= patchSize {
			continue
		}
		bigpatch := patches[x]
		// Remove the big old patch.
		patches = append(patches[:x], patches[x+1:]...)
		x--

		Start1 := bigpatch.Start1
		Start2 := bigpatch.Start2
		precontext := ""
		for len(bigpatch.diffs) != 0 {
			// Create one of several smaller patches.
			patch := Patch{}
			empty := true
			patch.Start1 = Start1 - len(precontext)
			patch.Start2 = Start2 - len(precontext)
			if len(precontext) != 0 {
				patch.Length1 = len(precontext)
				patch.Length2 = len(precontext)
				patch.diffs = append(patch.diffs, Diff{DiffEqual, precontext})
			}
			for len(bigpatch.diffs) != 0 && patch.Length1 < patchSize-dmp.PatchMargin {
				diffType := bigpatch.diffs[0].Type
				diffText := bigpatch.diffs[0].Text
				if diffType == DiffInsert {
					// Insertions are harmless.
					patch.Length2 += len(diffText)
					Start2 += len(diffText)
					patch.diffs = append(patch.diffs, bigpatch.diffs[0])
					bigpatch.diffs = bigpatch.diffs[1:]
					empty = false
				} else if diffType == DiffDelete && len(patch.diffs) == 1 && patch.diffs[0].Type == DiffEqual && len(diffText) > 2*patchSize {
					// This is a large deletion.  Let it pass in one chunk.
					patch.Length1 += len(diffText)
					Start1 += len(diffText)
					empty = false
					patch.diffs = append(patch.diffs, Diff{diffType, diffText})
					bigpatch.diffs = bigpatch.diffs[1:]
				} else {
					// Deletion or equality.  Only take as much as we can stomach.
					// rune-safety: align cutAt to nearest rune boundary so we never
					// slice diffText through a multi-byte character (the resulting
					// invalid UTF-8 is the root cause of issue #132 panic).
					cutAt := min(len(diffText), patchSize-patch.Length1-dmp.PatchMargin)
					cutAt = alignRuneStart(diffText, cutAt)
					// Loop-progress guard: when alignment collapses cutAt to 0 but
					// diffText is non-empty, advance by exactly one rune so the
					// outer for-loop is guaranteed to consume input each iteration.
					if cutAt == 0 && len(diffText) > 0 {
						_, size := utf8.DecodeRuneInString(diffText)
						cutAt = size
					}
					diffText = diffText[:cutAt]

					patch.Length1 += len(diffText)
					Start1 += len(diffText)
					if diffType == DiffEqual {
						patch.Length2 += len(diffText)
						Start2 += len(diffText)
					} else {
						empty = false
					}
					patch.diffs = append(patch.diffs, Diff{diffType, diffText})
					if diffText == bigpatch.diffs[0].Text {
						bigpatch.diffs = bigpatch.diffs[1:]
					} else {
						// rune-safety: len(diffText) was already rune-aligned above,
						// but align once more defensively for cutAt == 0 etc.
						remStart := alignRuneStart(bigpatch.diffs[0].Text, len(diffText))
						bigpatch.diffs[0].Text = bigpatch.diffs[0].Text[remStart:]
					}
				}
			}
			// Compute the head context for the next patch.
			precontext = dmp.DiffText2(patch.diffs)
			// rune-safety: align precontext start to nearest rune boundary.
			precontext = precontext[alignRuneStart(precontext, max(0, len(precontext)-dmp.PatchMargin)):]

			// Append the end context for this patch.
			postcontext := diffText1Prefix(bigpatch.diffs, dmp.PatchMargin)

			if len(postcontext) != 0 {
				patch.Length1 += len(postcontext)
				patch.Length2 += len(postcontext)
				if len(patch.diffs) != 0 && patch.diffs[len(patch.diffs)-1].Type == DiffEqual {
					patch.diffs[len(patch.diffs)-1].Text += postcontext
				} else {
					patch.diffs = append(patch.diffs, Diff{DiffEqual, postcontext})
				}
			}
			if !empty {
				x++
				patches = append(patches[:x], append([]Patch{patch}, patches[x:]...)...)
			}
		}
	}
	return patches
}

// diffText1Prefix returns the first n bytes of DiffText1(diffs) without building the whole text.
// If byte n falls inside a multi-byte rune, the prefix stops before that rune.
func diffText1Prefix(diffs []Diff, n int) string {
	var text strings.Builder
	for _, aDiff := range diffs {
		if aDiff.Type == DiffInsert {
			continue
		}
		remaining := n - text.Len()
		if remaining <= 0 {
			break
		}
		if len(aDiff.Text) >= remaining {
			cut := alignRuneStart(aDiff.Text, remaining)
			if text.Len() == 0 {
				return aDiff.Text[:cut]
			}
			_, _ = text.WriteString(aDiff.Text[:cut])
			break
		}
		_, _ = text.WriteString(aDiff.Text)
	}
	return text.String()
}

// PatchToText takes a list of patches and returns a textual representation.
func (dmp *DiffMatchPatch) PatchToText(patches []Patch) string {
	var text bytes.Buffer
	for _, aPatch := range patches {
		_, _ = text.WriteString(aPatch.String())
	}
	return text.String()
}

var patchHeaderRegex = regexp.MustCompile("^@@ -(\\d+),?(\\d*) \\+(\\d+),?(\\d*) @@$")

// PatchFromText parses a textual representation of patches and returns a List of Patch objects.
func (dmp *DiffMatchPatch) PatchFromText(textline string) ([]Patch, error) {
	patches := []Patch{}
	if len(textline) == 0 {
		return patches, nil
	}
	text := strings.Split(textline, "\n")
	textPointer := 0

	var patch Patch
	var sign uint8
	var line string
	for textPointer < len(text) {
		m := patchHeaderRegex.FindStringSubmatch(text[textPointer])
		if m == nil {
			return patches, errors.New("Invalid patch string: " + text[textPointer])
		}

		patch = Patch{}

		// Coordinates are parsed as 32-bit integers, like the reference Java
		// and C# implementations. Out-of-range values are invalid rather
		// than clamped: accepting them made hostile patches apply at wrong
		// locations and leak internal padding into the result.
		start1, err := strconv.ParseInt(m[1], 10, 32)
		if err != nil {
			return patches, errors.New("Invalid patch string: " + text[textPointer])
		}
		patch.Start1 = int(start1)
		if len(m[2]) == 0 {
			patch.Start1--
			patch.Length1 = 1
		} else if m[2] == "0" {
			patch.Length1 = 0
		} else {
			length1, err := strconv.ParseInt(m[2], 10, 32)
			if err != nil {
				return patches, errors.New("Invalid patch string: " + text[textPointer])
			}
			patch.Start1--
			patch.Length1 = int(length1)
		}

		start2, err := strconv.ParseInt(m[3], 10, 32)
		if err != nil {
			return patches, errors.New("Invalid patch string: " + text[textPointer])
		}
		patch.Start2 = int(start2)
		if len(m[4]) == 0 {
			patch.Start2--
			patch.Length2 = 1
		} else if m[4] == "0" {
			patch.Length2 = 0
		} else {
			length2, err := strconv.ParseInt(m[4], 10, 32)
			if err != nil {
				return patches, errors.New("Invalid patch string: " + text[textPointer])
			}
			patch.Start2--
			patch.Length2 = int(length2)
		}
		textPointer++

		for textPointer < len(text) {
			if len(text[textPointer]) > 0 {
				sign = text[textPointer][0]
			} else {
				textPointer++
				continue
			}

			line = text[textPointer][1:]
			line = strings.Replace(line, "+", "%2b", -1)
			var err error
			line, err = url.QueryUnescape(line)
			if err != nil {
				return patches, err
			}
			if sign == '-' {
				// Deletion.
				patch.diffs = append(patch.diffs, Diff{DiffDelete, line})
			} else if sign == '+' {
				// Insertion.
				patch.diffs = append(patch.diffs, Diff{DiffInsert, line})
			} else if sign == ' ' {
				// Minor equality.
				patch.diffs = append(patch.diffs, Diff{DiffEqual, line})
			} else if sign == '@' {
				// Start of next patch.
				break
			} else {
				// WTF?
				return patches, errors.New("Invalid patch mode '" + string(sign) + "' in: " + string(line))
			}
			textPointer++
		}

		patches = append(patches, patch)
	}
	return patches, nil
}
