package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLocateSection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveHaystack []string
		giveNeedle   []string
		name         string
		want         Option[int]
	}{
		{
			giveHaystack: []string{
				"one",
				"two",
				"Some text",
			},
			giveNeedle: []string{
				"one",
				"two",
			},
			name: "section at beginning",
			want: Some(0),
		},
		{
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
			name: "section in middle",
			want: Some(1),
		},
		{
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
			name: "section at end",
			want: Some(2),
		},
		{
			giveHaystack: []string{
				"one",
				"three",
			},
			giveNeedle: []string{
				"one",
				"two",
			},
			name: "section not found",
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
