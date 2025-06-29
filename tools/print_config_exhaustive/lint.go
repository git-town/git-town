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
	normalConfigDataPath = "internal/config/configdomain/normal_config_data.go"
	printConfigPath      = "internal/cmd/config/root.go"
)

var whiteList = []string{
	"Aliases",
	"BranchTypeOverrides",
}

func main() {
	// read and parse definition file
	definitionText := readFile(normalConfigDataPath)
	definitionFields := DefinitionFields(definitionText)

	// read and parse print file
	printText := readFile(printConfigPath)
	printBody := ParsePrintFile(printText)

	// find fields that are in the definition file but not the print file
	unprinted := FindUnprinted(definitionFields, printBody, whiteList)
	if len(unprinted) > 0 {
		printMissingFields(unprinted)
		os.Exit(1)
	}
}

func DefinitionFields(text string) []string {
	structRE := regexp.MustCompile(`type NormalConfigData struct {([^}]*)}`)
	match := structRE.FindStringSubmatch(text)
	if len(match) < 2 {
		log.Fatalf("Error: Failed to find NormalConfigData struct")
	}
	result := []string{}
	for _, line := range strings.Split(strings.TrimSpace(match[1]), "\n") {
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

func FindUnprinted(fields []string, text string, whiteList []string) []string {
	result := []string{}
	for _, field := range fields {
		if strings.Contains(text, field) {
			continue
		}
		if isWhitelisted(field, whiteList) {
			continue
		}
		result = append(result, field)
	}
	return result
}

func ParsePrintFile(text string) string {
	functionContentRE := regexp.MustCompile(`func printConfig\(.*?\) {([^}]*)}`)
	match := functionContentRE.FindStringSubmatch(text)
	if len(match) < 2 {
		log.Fatalf("Error: Failed to find printConfig function")
	}
	return match[1]
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
