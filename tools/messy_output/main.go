package main

import "fmt"

func main() {
	// find all .feature files
	// - has tag @messyoutput
	// find all scenarios in the feature file
	// - has tag @messyoutput
	// - has step step `I run "git-town append new" and enter into the dialog`
	// list all scenarios that have the step but no tag
	// list all scenarios that have the tag but no step
	fmt.Println("hello")
}
