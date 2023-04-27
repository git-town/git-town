package datatable

import "github.com/git-town/git-town/v8/test/commands"

// FilesInBranches provides a data table of files and their content in all branches.
func FilesInBranches(cmds commands.Repo, mainBranch string) (DataTable, error) {
	result := DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := cmds.ProdGit().LocalBranchesMainFirst(mainBranch)
	if err != nil {
		return DataTable{}, err
	}
	lastBranch := ""
	for _, branch := range branches {
		files, err := commands.FilesInBranch(cmds, branch)
		if err != nil {
			return DataTable{}, err
		}
		for _, file := range files {
			content, err := commands.FileContentInCommit(cmds, branch, file)
			if err != nil {
				return DataTable{}, err
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
