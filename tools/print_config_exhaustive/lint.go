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
	definitionFields := definitionFields(definitionText)

	// read and parse print file
	printText := readFile(printConfigPath)
	printBody := parsePrintFile(printText)

	// find fields that are in the definition file but not the print file
	unprinted := findUnprinted(definitionFields, printBody, whiteList)
	if len(unprinted) == 0 {
		return
	}
	fmt.Println("Missing fields in printConfig function:")
	for _, field := range unprinted {
		fmt.Println(field)
	}
	os.Exit(1)
}

func readFile(path string) string {
	result, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading %s: %v\n", path, err)
	}
	return string(result)
}

func definitionFields(text string) []string {
	structRE := regexp.MustCompile(`type NormalConfigData struct {([^}]*)}`)
	match := structRE.FindStringSubmatch(string(text))
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

func parsePrintFile(text string) string {
	functionContentRE := regexp.MustCompile(`func printConfig\(.*?\) {([^}]*)}`)
	match := functionContentRE.FindStringSubmatch(string(text))
	if len(match) < 2 {
		log.Fatalf("Error: Failed to find printConfig function")
	}
	return match[1]
}

func findUnprinted(fields []string, text string, whiteList []string) []string {
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

func isWhitelisted(field string, whitelist []string) bool {
	return slices.Contains(whitelist, field)
}
