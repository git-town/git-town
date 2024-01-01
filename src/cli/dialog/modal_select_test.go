package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/stretchr/testify/assert"
)

func TestModalEntries(t *testing.T) {
	t.Run("Texts", func(t *testing.T) {
		mes := dialog.ModalEntries{
			dialog.ModalEntry{
				Text: "One",
				Value: "one",
			},
			dialog.ModalEntry{
				Text: "Two",
				Value: "two",
			},
		}
		want := []string{"One", "Two"}
		have := mes.Texts()
		assert.Equal(t, want, have)
	})
}
