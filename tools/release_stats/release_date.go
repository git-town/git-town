package main

import (
	"os/exec"
	"strings"
	"time"
)

func releaseDate(tag string) time.Time {
	cmd := exec.Command("git", "log", "-1", "--format=%cI", tag)
	outputData, err := cmd.CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	output := strings.TrimSpace(string(outputData))
	result, err := time.Parse(time.RFC3339, output)
	if err != nil {
		panic(err.Error())
	}
	return result
}
