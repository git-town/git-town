package commands

import "github.com/git-town/git-town/v8/test/datatable"

// FilesInBranches provides a data table of files and their content in all branches.
func FilesInBranches(cmds *TestCommands, mainBranch string) (datatable.DataTable, error) {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := cmds.LocalBranchesMainFirst(mainBranch)
	if err != nil {
		return datatable.DataTable{}, err
	}
	lastBranch := ""
	for _, branch := range branches {
		files, err := FilesInBranch(cmds, branch)
		if err != nil {
			return datatable.DataTable{}, err
		}
		for _, file := range files {
			content, err := FileContentInCommit(cmds, branch, file)
			if err != nil {
				return datatable.DataTable{}, err
			}
			if branch == lastBranch {
				result.AddRow("", file, content)
			} else {
				result.AddRow(branch, file, content)
			}
			lastBranch = branch
		}
	}
	return result, err
}
