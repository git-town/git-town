package git

import (
	"fmt"
	"log"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/test/helpers"
)

// Commit describes a Git commit.
type Commit struct {
	Author      Option[string] `exhaustruct:"optional"`
	Branch      gitdomain.LocalBranchName
	FileContent Option[string]    `exhaustruct:"optional"`
	FileName    Option[string]    `exhaustruct:"optional"`
	Locations   Option[Locations] `exhaustruct:"optional"`
	Message     string
	SHA         Option[gitdomain.SHA] `exhaustruct:"optional"`
}

func (self Commit) GetAuthor() string {
	if author, hasAuthor := self.Author.Get(); hasAuthor {
		return author
	}
	return "user <email@example.com>"
}

func (self Commit) GetFileContent() string {
	if content, hasContent := self.FileContent.Get(); hasContent {
		return content
	}
	return fmt.Sprintf("default file content for file %q in branch %q", self.GetFileName(), self.Branch)
}

func (self Commit) GetFileName() string {
	if name, hasName := self.FileName.Get(); hasName {
		return name
	}
	return fmt.Sprintf("default_filename_%s_%s", self.Branch, self.GetLocations().Join("_"))
}

func (self Commit) GetLocations() Locations {
	if locations, hasLocations := self.Locations.Get(); hasLocations {
		return locations
	}
	return Locations{LocationLocal, LocationOrigin}
}

func (self Commit) GetSHA() gitdomain.SHA {
	if sha, hasSHA := self.SHA.Get(); hasSHA {
		return sha
	}
	return gitdomain.NewSHA("111111")
}

func (self *Commit) SetFileName(fileName string) {
	self.FileName = Some(fileName)
}

// FromGherkinTable provides a Commit collection representing the data in the given Gherkin table.
func FromGherkinTable(table *godog.Table, branchName gitdomain.LocalBranchName) []Commit {
	columnNames := helpers.TableFields(table)
	lastBranch := ""
	lastLocationName := ""
	result := []Commit{}
	for _, row := range table.Rows[1:] {
		commit := Commit{
			Author:      None[string](),
			Branch:      branchName,
			FileContent: None[string](),
			FileName:    None[string](),
			Locations:   None[Locations](),
			Message:     "default commit message",
			SHA:         None[gitdomain.SHA](),
		}
		for cellNo, cell := range row.Cells {
			columnName := columnNames[cellNo]
			cellValue := cell.Value
			switch columnName {
			case "BRANCH":
				if cellValue == "" {
					cellValue = lastBranch
				} else {
					lastBranch = cellValue
				}
				commit.Branch = gitdomain.NewLocalBranchName(cellValue)
			case "LOCATION":
				if cell.Value == "" {
					cellValue = lastLocationName
				} else {
					lastLocationName = cellValue
				}
				commit.Locations = Some(NewLocations(cellValue))
			case "MESSAGE":
				commit.Message = cellValue
			case "FILE NAME":
				commit.FileName = Some(cellValue)
			case "FILE CONTENT":
				commit.FileContent = Some(cellValue)
			case "AUTHOR":
				commit.Author = Some(cellValue)
			default:
				log.Fatalf("unknown Commit property: %s", columnName)
			}
		}
		result = append(result, commit)
	}
	return result
}
