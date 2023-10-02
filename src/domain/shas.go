package domain

import "strings"

type SHAs []SHA

func (ss SHAs) Join(sep string) string {
	return strings.Join(ss.Strings(), sep)
}

func (ss SHAs) Strings() []string {
	result := make([]string, len(ss))
	for s, sha := range ss {
		result[s] = sha.String()
	}
	return result
}
