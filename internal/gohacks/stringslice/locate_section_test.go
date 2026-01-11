package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLocateSection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		giveHaystack []string
		giveNeedle   []string
		want         Option[int]
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
			want: Some(0),
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
			want: Some(1),
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
			want: Some(2),
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
			want: None[int](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have := stringslice.LocateSection(tt.giveHaystack, tt.giveNeedle)
			must.True(t, tt.want.Equal(have))
		})
	}
}
