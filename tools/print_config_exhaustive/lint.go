package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

const (
	normalConfigDataPath = "internal/config/normal_config.go"
	printConfigPath      = "internal/cmd/config/root.go"
)

var whiteList = []string{
	"Aliases",
	"BranchTypeOverrides",
	"DryRun",
	"ProposalsShowLineage", // TODO: remove once https://github.com/git-town/git-town/issues/3003 is shipped
	"Verbose",
}

func main() {
	definedFields := FindDefinedFields(readFile(normalConfigDataPath))
	printFunction := FindPrintFunc(readFile(printConfigPath))
	unprinted := FindUnprintedFields(definedFields, printFunction, whiteList)
	if len(unprinted) > 0 {
		printMissingFields(unprinted)
		os.Exit(1)
	}
}

func FindDefinedFields(text string) []string {
	structRE := regexp.MustCompile(`type NormalConfig struct {([^}]*)}`)
	match := structRE.FindStringSubmatch(text)
	if len(match) < 2 {
		log.Fatalf("Error: Failed to find NormalConfig struct")
	}
	result := []string{}
	for line := range strings.SplitSeq(strings.TrimSpace(match[1]), "\n") {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "//") {
			continue // skip comments
		}
		parts := strings.Fields(trimmedLine)
		if len(parts) > 0 {
			result = append(result, parts[0])
		}
	}
	return result
}

func FindPrintFunc(text string) string {
	funcBodyRE := regexp.MustCompile(`func printConfig\(.*?\) {([^}]*)}`)
	match := funcBodyRE.FindStringSubmatch(text)
	if len(match) < 2 {
		log.Fatalf("Error: Failed to find printConfig function")
	}
	return match[1]
}

func FindUnprintedFields(fields []string, text string, whiteList []string) []string {
	result := []string{}
	for _, field := range fields {
		if !isPrinted(field, text) && !isWhitelisted(field, whiteList) {
			result = append(result, field)
		}
	}
	return result
}

func isPrinted(field, text string) bool {
	return strings.Contains(text, field)
}

func isWhitelisted(field string, whitelist []string) bool {
	return slices.Contains(whitelist, field)
}

func printMissingFields(unprinted []string) {
	fmt.Println("Missing fields in printConfig function:")
	for _, field := range unprinted {
		fmt.Println(field)
	}
}

func readFile(path string) string {
	result, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading %s: %v\n", path, err)
	}
	return string(result)
}
