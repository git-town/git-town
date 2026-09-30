package subshell

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestEscapeCmdMetaChars(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		give string
		want string
	}{
		"no metacharacters": {
			give: "https://example.com/path",
			want: "https://example.com/path",
		},
		"ampersand in query string": {
			give: "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Amain",
			want: "https://bitbucket.org/org/repo/pull-requests/new?source=branch^&dest=org%2Frepo%3Amain",
		},
		"multiple ampersands": {
			give: "https://example.com?a=1&b=2&c=3",
			want: "https://example.com?a=1^&b=2^&c=3",
		},
		"caret is escaped first": {
			give: "foo^bar&baz",
			want: "foo^^bar^&baz",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			must.Eq(t, tc.want, escapeCmdMetaChars(tc.give))
		})
	}
}
