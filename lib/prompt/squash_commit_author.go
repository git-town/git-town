package prompt

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/util"
	"github.com/fatih/color"
)

func GetSquashCommitAuthor(branchName string) string {
	authors := getBranchAuthors(branchName)
	if len(authors) == 1 {
		return authors[0].NameAndEmail
	} else {
		fmt.Printf(squashCommitAuthorHeaderTemplate, branchName)
		printNumberedAuthors(authors)
		fmt.Println()
		return askForAuthor(authors)
	}
}

// Helpers

type Author struct {
	NameAndEmail    string
	NumberOfCommits string
}

var squashCommitAuthorHeaderTemplate = `
Multiple people authored the '%s' branch.
Please choose an author for the squash commit.

`

func askForAuthor(authors []Author) string {
	for {
		fmt.Print("Enter user's number or a custom author (default: 1): ")
		author := parseAuthor(util.GetUserInput(), authors)
		if author != "" {
			return author
		}
	}
}

func getBranchAuthors(branchName string) (result []Author) {
	output := util.GetCommandOutput("git", "shortlog", "-s", "-n", "-e", git.GetMainBranch()+".."+branchName)
	for _, line := range strings.Split(output, "\n") {
		line := strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, Author{NameAndEmail: parts[1], NumberOfCommits: parts[0]})
	}
	return
}

func parseAuthor(userInput string, authors []Author) string {
	numericRegex, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		log.Fatal("Error compiling numeric regular expression: ", err)
	}

	if numericRegex.MatchString(userInput) {
		return parseAuthorNumber(userInput, authors)
	} else if userInput == "" {
		return authors[0].NameAndEmail
	} else {
		return userInput
	}
}

func parseAuthorNumber(userInput string, authors []Author) string {
	index, err := strconv.Atoi(userInput)
	if err != nil {
		log.Fatal("Error parsing string to integer: ", err)
	}
	if index >= 1 && index <= len(authors) {
		return authors[index-1].NameAndEmail
	} else {
		util.PrintError("Invalid author number")
		return ""
	}
}

func printNumberedAuthors(authors []Author) {
	boldFmt := color.New(color.Bold)
	for index, author := range authors {
		stat := util.Pluralize(author.NumberOfCommits, "commit")
		fmt.Printf("  %s: %s (%s)\n", boldFmt.Sprintf("%d", index+1), author.NameAndEmail, stat)
	}
}
