package gitdomain

import "strings"

type SHAs []SHA

func (self SHAs) First() SHA {
	return self[0]
}

func (self SHAs) Join(sep string) string {
	return strings.Join(self.Strings(), sep)
}

func (self SHAs) Last() SHA {
	return self[len(self)]
}

func (self SHAs) Strings() []string {
	result := make([]string, len(self))
	for s, sha := range self {
		result[s] = sha.String()
	}
	return result
}
