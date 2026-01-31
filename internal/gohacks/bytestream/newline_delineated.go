package bytestream

import "bytes"

// NewlineDelineated is raw Git output that is delineated by newlines.
type NewlineDelineated []byte

func (self NewlineDelineated) Sanitize() Sanitized {
	lines := bytes.Split(self, []byte("\n"))
	secretKeys := [][]byte{
		[]byte("git-town.github-token"),
		[]byte("git-town.gitlab-token"),
		[]byte("git-town.forgejo-token"),
		[]byte("git-town.bitbucket-app-password"),
		[]byte("git-town.gitea-token"),
		[]byte("user.email"),
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
