package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestUseCustomMessageVerbatim(t *testing.T) {
	t.Parallel()
	want := gitdomain.CommitMessage("foo")
	useMessage := configdomain.UseCustomMessageVerbatim(want)
	message, isCustom := useMessage.GetCustomMessage()
	must.True(t, isCustom)
	must.Eq(t, want, message)
	must.True(t, useMessage.IsCleanupVerbatim())
	must.True(t, useMessage.IsCustomMessage())
	must.False(t, useMessage.IsEditDefault())
	must.False(t, useMessage.IsUseDefault())
}

func TestUseMessageCleanupVerbatimDefaultsToFalse(t *testing.T) {
	t.Parallel()
	useMessage := configdomain.UseCustomMessage("foo")
	must.False(t, useMessage.IsCleanupVerbatim())
}
