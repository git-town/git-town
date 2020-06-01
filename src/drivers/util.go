package drivers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func printLog(message string) {
	fmt.Println()
	_, err := color.New(color.Bold).Println(message)
	if err != nil {
		panic(err)
	}
}

// GetURLHostname returns the hostname contained within the given Git URL.
func GetURLHostname(url string) string {
	hostnameRegex := regexp.MustCompile("(^[^:]*://([^@]*@)?|git@)([^/:]+).*")
	matches := hostnameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return matches[3]
}

// GetURLRepositoryName returns the repository name contains within the given Git URL.
func GetURLRepositoryName(url string) string {
	hostname := GetURLHostname(url)
	repositoryNameRegex := regexp.MustCompile(".*" + hostname + "[/:](.+)")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}