package domain

type RemoteBranchNames []RemoteBranchName

// Strings provides these remote branch names as strings.
func (r RemoteBranchNames) Strings() []string {
	result := make([]string, len(r))
	for b, branch := range r {
		result[b] = branch.String()
	}
	return result
}
