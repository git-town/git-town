package domain

import "strings"

type SHAs []SHA

func (self SHAs) Join(sep string) string {
	return strings.Join(self.Strings(), sep)
}

func (self SHAs) Strings() []string {
	result := make([]string, len(self))
	for s, sha := range self {
		result[s] = sha.String()
	}
	return result
}
