package domain

import "strings"

type SHAs []SHA

func (shas SHAs) Join(sep string) string {
	return strings.Join(shas.Strings(), sep)
}

func (shas SHAs) Strings() []string {
	result := make([]string, len(shas))
	for s, sha := range shas {
		result[s] = sha.String()
	}
	return result
}
