package gitdomain

import "strings"

type SHAs []SHA

func NewSHAs(ids ...string) SHAs {
	result := make(SHAs, len(ids))
	for i, id := range ids {
		result[i] = NewSHA(id)
	}
	return result
}

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
