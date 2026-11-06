package envconfig_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v24/internal/config/envconfig"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/shoenig/test/must"
)

func TestGhPrompt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		env  []string
		want forgedomain.GhPrompt
	}{
		{env: []string{}, want: true},
		{env: []string{"OTHER=1"}, want: true},
		{env: []string{"GH_PROMPT_DISABLED="}, want: false},
		{env: []string{"GH_PROMPT_DISABLED=0"}, want: false},
		{env: []string{"GH_PROMPT_DISABLED=1"}, want: false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.env), func(t *testing.T) {
			t.Parallel()
			have := envconfig.GhPrompt(envconfig.NewEnvVars(test.env))
			must.EqOp(t, test.want, have)
		})
	}
}
