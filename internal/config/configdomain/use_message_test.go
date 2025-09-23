package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestUseMessage(t *testing.T) {
	t.Parallel()

	t.Run("EditDefaultMessage", func(t *testing.T) {
		t.Parallel()

		m := configdomain.EditDefaultMessage()

		// Custom type getters.
		_, isCustom := m.GetCustomMessage()
		must.False(t, isCustom)
		// Type checks.
		must.False(t, m.IsCustomMessage())
		must.True(t, m.IsEditDefault())
		must.False(t, m.IsUseDefault())
	})

	t.Run("GetCustomMessageOrPanic", func(t *testing.T) {
		// Arrange.
		t.Parallel()
		m := configdomain.EditDefaultMessage()

		// Assert.
		defer func() {
			r := recover()
			must.Eq(t, r, "UseMessage is not UseCustomMessage")
		}()

		// Act.
		_ = m.GetCustomMessageOrPanic()
	})

	t.Run("UseCustomMessage", func(t *testing.T) {
		t.Parallel()
		want := gitdomain.CommitMessage("foo")

		m := configdomain.UseCustomMessage(want)

		// Custom type getters.
		message, isCustom := m.GetCustomMessage()
		must.True(t, isCustom)
		must.Eq(t, want, message)
		must.Eq(t, want, m.GetCustomMessageOrPanic())
		// Type checks.
		must.True(t, m.IsCustomMessage())
		must.False(t, m.IsEditDefault())
		must.False(t, m.IsUseDefault())
	})

	t.Run("UseCustomMessageOr", func(t *testing.T) {
		t.Parallel()
		t.Run("WithSome", func(t *testing.T) {
			t.Parallel()
			want := gitdomain.CommitMessage("foo")

			m := configdomain.UseCustomMessageOr(Some(want), configdomain.EditDefaultMessage())

			message, isCustom := m.GetCustomMessage()
			must.True(t, isCustom)
			must.Eq(t, want, message)
		})
		t.Run("WithNone", func(t *testing.T) {
			t.Parallel()

			m := configdomain.UseCustomMessageOr(None[gitdomain.CommitMessage](), configdomain.EditDefaultMessage())

			must.True(t, m.IsEditDefault())
		})
	})

	t.Run("UseDefaultMessage", func(t *testing.T) {
		t.Parallel()

		m := configdomain.UseDefaultMessage()

		// Custom type getters.
		_, isCustom := m.GetCustomMessage()
		must.False(t, isCustom)
		// Type checks.
		must.False(t, m.IsCustomMessage())
		must.False(t, m.IsEditDefault())
		must.True(t, m.IsUseDefault())
	})

	t.Run("UseMessageWithFallbackToDefault", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			desc              string
			message           Option[gitdomain.CommitMessage]
			fallbackToDefault bool
			wantIsCustom      bool
			wantIsEditDefault bool
			wantIsUseDefault  bool
		}{
			{
				desc:              "No message, no fallback",
				message:           None[gitdomain.CommitMessage](),
				fallbackToDefault: false,
				wantIsCustom:      false,
				wantIsEditDefault: true,
				wantIsUseDefault:  false,
			}, {
				desc:              "No message, with fallback",
				message:           None[gitdomain.CommitMessage](),
				fallbackToDefault: true,
				wantIsCustom:      false,
				wantIsEditDefault: false,
				wantIsUseDefault:  true,
			}, {
				desc:              "With message, no fallback",
				message:           Some(gitdomain.CommitMessage("foo")),
				fallbackToDefault: false,
				wantIsCustom:      true,
				wantIsEditDefault: false,
				wantIsUseDefault:  false,
			}, {
				desc:              "With message, with fallback",
				message:           Some(gitdomain.CommitMessage("foo")),
				fallbackToDefault: true,
				wantIsCustom:      true,
				wantIsEditDefault: false,
				wantIsUseDefault:  false,
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				t.Parallel()
				m := configdomain.UseMessageWithFallbackToDefault(tc.message, tc.fallbackToDefault)

				message, isCustom := m.GetCustomMessage()
				must.Eq(t, tc.wantIsCustom, isCustom)
				if isCustom {
					must.Eq(t, tc.message.GetOrPanic(), message)
					must.Eq(t, tc.message.GetOrPanic(), m.GetCustomMessageOrPanic())
				}
				must.Eq(t, tc.wantIsCustom, m.IsCustomMessage())
				must.Eq(t, tc.wantIsEditDefault, m.IsEditDefault())
				must.Eq(t, tc.wantIsUseDefault, m.IsUseDefault())
			})
		}
	})
}
