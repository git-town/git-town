package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestLocateSection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveHaystack []string
		giveNeedle   []string
		wantIdx      int
		wantFound    bool
	}{
		{
			name: "section at beginning",
			giveHaystack: []string{
				"one",
				"two",
				"Some text",
			},
			giveNeedle: []string{
				"one",
				"two",
			},
			wantIdx:   0,
			wantFound: true,
		},
		{
			name: "section in middle",
			giveHaystack: []string{
				"Some text",
				"    one",
				"    two",
				"More text",
			},
			giveNeedle: []string{
				"one",
				"two",
			},
			wantIdx:   1,
			wantFound: true,
		},
		{
			name: "section at end",
			giveHaystack: []string{
				"Some text",
				"More text",
				"  one",
				"  two",
			},
			giveNeedle: []string{
				"one",
				"two",
			},
			wantIdx:   2,
			wantFound: true,
		},
		{
			name: "section not found",
			giveHaystack: []string{
				"one",
				"three",
			},
			giveNeedle: []string{
				"one",
				"two",
			},
			wantIdx:   -1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			haveIdx, haveFound := stringslice.LocateSection(tt.giveHaystack, tt.giveNeedle)
			must.EqOp(t, tt.wantIdx, haveIdx)
			must.EqOp(t, tt.wantFound, haveFound)
		})
	}
}
