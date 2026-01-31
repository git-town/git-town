package bytestream

import (
	"bytes"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

// NewlineDelineated is raw Git output that is delineated by newlines.
type NewlineDelineated []byte

func (self NewlineDelineated) Sanitize() Sanitized {
	lines := bytes.Split(self, []byte("\n"))
	secretKeys := [][]byte{
		[]byte(configdomain.KeyBitbucketAppPassword),
		[]byte(configdomain.KeyDeprecatedCodebergToken),
		[]byte(configdomain.KeyForgejoToken),
		[]byte(configdomain.KeyGiteaToken),
		[]byte(configdomain.KeyGithubToken),
		[]byte(configdomain.KeyGitlabToken),
		[]byte(configdomain.KeyGitUserEmail),
	}
	for i, line := range lines {
		for _, key := range secretKeys {
			if bytes.Equal(line, key) && i+1 < len(lines) {
				lines[i+1] = []byte("(redacted)")
				break
			}
		}
	}
	return bytes.Join(lines, []byte("\n"))
}
