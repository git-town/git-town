package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/util"
)

// Runner executes Git commands.
type Runner struct {
	command.Shell                        // for running console commands
	remoteBranchCache *RemoteBranchCache // caches the remote branches of this Git repo
	*Configuration                       // caches Git configuration settings
}

// NewRunner provides Runner instances.
func NewRunner(shell command.Shell, remoteBranches *RemoteBranchCache, config *Configuration) Runner {
	return Runner{shell, remoteBranches, config}
}

// AbortMerge cancels a currently ongoing Git merge operation.
func (r *Runner) AbortMerge() error {
	res, err := r.Run("git", "merge", "--abort")
	if err != nil {
		return fmt.Errorf("cannot abort current merge: %w\n%s", err, res.Output())
	}
	return nil
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (r *Runner) AbortRebase() error {
	res, err := r.Run("git", "rebase", "--abort")
	if err != nil {
		return fmt.Errorf("cannot abort current merge: %w\n%s", err, res.Output())
	}
	return nil
}

// AddRemote adds the given Git remote to this repository.
func (r *Runner) AddRemote(name, value string) error {
	res, err := r.Run("git", "remote", "add", name, value)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w\n%s", name, value, err, res.Output())
	}
	return nil
}

// BranchHasUnmergedCommits indicates whether the branch with the given name
// contains commits that are not merged into the main branch
func (r *Runner) BranchHasUnmergedCommits(branch string) (bool, error) {
	out, err := r.Run("git", "log", r.GetMainBranch()+".."+branch)
	if err != nil {
		return false, fmt.Errorf("cannot determine if branch %q has unmerged commits: %w\n%s", branch, err, out.Output())
	}
	return out.OutputSanitized() != "", nil
}

// BranchSha provides the SHA for the local branch with the given name.
func (r *Runner) BranchSha(name string) (sha string, err error) {
	outcome, err := r.Run("git", "rev-parse", name)
	if err != nil {
		return "", fmt.Errorf("cannot determine SHA of local branch %q: %w\n%s", name, err, outcome.Output())
	}
	return outcome.OutputSanitized(), nil
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (r *Runner) CheckoutBranch(name string) error {
	outcome, err := r.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q in repo %q: %w\n%v", name, r.WorkingDir(), err, outcome)
	}
	currentBranchCache = name
	return nil
}

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided
func (r *Runner) CommentOutSquashCommitMessage(prefix string) error {
	squashMessageFile := ".git/SQUASH_MSG"
	contentBytes, err := ioutil.ReadFile(squashMessageFile)
	if err != nil {
		return fmt.Errorf("cannot read squash message file %q: %w", squashMessageFile, err)
	}
	content := string(contentBytes)
	if prefix != "" {
		content = prefix + "\n" + content
	}
	content = regexp.MustCompile("(?m)^").ReplaceAllString(content, "# ")
	return ioutil.WriteFile(squashMessageFile, []byte(content), 0644)
}

// CommitNoEdit commits all staged files with the default commit message.
func (r *Runner) CommitNoEdit() error {
	outcome, err := r.Run("git", "commit", "--no-edit")
	if err != nil {
		return fmt.Errorf("cannot commit files: %w\n%s", err, outcome.Output())
	}
	return nil
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (r *Runner) Commits(fields []string) (result []Commit, err error) {
	branches, err := r.LocalBranches()
	if err != nil {
		return result, fmt.Errorf("cannot determine the Git branches: %w", err)
	}
	for _, branch := range branches {
		commits, err := r.CommitsInBranch(branch, fields)
		if err != nil {
			return result, err
		}
		result = append(result, commits...)
	}
	return result, nil
}

// CommitsInBranch provides all commits in the given Git branch.
func (r *Runner) CommitsInBranch(branch string, fields []string) (result []Commit, err error) {
	outcome, err := r.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return result, fmt.Errorf("cannot get commits in branch %q: %w", branch, err)
	}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
		parts := strings.Split(line, "|")
		commit := Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if util.DoesStringArrayContain(fields, "FILE NAME") {
			filenames, err := r.FilesInCommit(commit.SHA)
			if err != nil {
				return result, fmt.Errorf("cannot determine file name for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileName = strings.Join(filenames, ", ")
		}
		if util.DoesStringArrayContain(fields, "FILE CONTENT") {
			filecontent, err := r.FileContentInCommit(commit.SHA, commit.FileName)
			if err != nil {
				return result, fmt.Errorf("cannot determine file content for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result, nil
}

// CommitStagedChanges commits the currently staged changes.
func (r *Runner) CommitStagedChanges(message string) error {
	var out *command.Result
	var err error
	if message != "" {
		out, err = r.Run("git", "commit", "-m", message)
	} else {
		out, err = r.Run("git", "commit", "--no-edit")
	}
	if err != nil {
		return fmt.Errorf("cannot commit staged changes: %w\n%s", err, out)
	}
	return nil
}

// ConnectTrackingBranch connects the branch with the given name to its remote tracking branch.
// The branch must exist.
func (r *Runner) ConnectTrackingBranch(name string) error {
	out, err := r.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	if err != nil {
		return fmt.Errorf("cannot connect tracking branch for %q: %w\n%s", name, err, out)
	}
	return nil
}

// ContinueRebase continues the currently ongoing rebase.
func (r *Runner) ContinueRebase() error {
	outcome, err := r.Run("git", "rebase", "--continue")
	if err != nil {
		return fmt.Errorf("cannot continue rebase: %w\n%s", err, outcome.Output())
	}
	return nil
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (r *Runner) CreateBranch(name, parent string) error {
	outcome, err := r.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (r *Runner) CreateChildFeatureBranch(name string, parent string) error {
	err := r.CreateBranch(name, parent)
	if err != nil {
		return fmt.Errorf("cannot create child branch %q: %w", name, err)
	}
	_ = r.Configuration.SetParentBranch(name, parent)
	return nil
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (r *Runner) CreateCommit(commit Commit) error {
	err := r.CheckoutBranch(commit.Branch)
	if err != nil {
		return fmt.Errorf("cannot checkout branch %q: %w", commit.Branch, err)
	}
	err = r.CreateFile(commit.FileName, commit.FileContent)
	if err != nil {
		return fmt.Errorf("cannot create file %q needed for commit: %w", commit.FileName, err)
	}
	outcome, err := r.Run("git", "add", commit.FileName)
	if err != nil {
		return fmt.Errorf("cannot add file to commit: %w\n%v", err, outcome)
	}
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	outcome, err = r.Run("git", commands...)
	if err != nil {
		return fmt.Errorf("cannot commit: %w\n%v", err, outcome)
	}
	return nil
}

// CreateFeatureBranch creates a feature branch with the given name in this repository.
func (r *Runner) CreateFeatureBranch(name string) error {
	err := r.RunMany([][]string{
		{"git", "branch", name, "main"},
		{"git", "config", "git-town-branch." + name + ".parent", "main"},
	})
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w", name, err)
	}
	return nil
}

// CreateFeatureBranchNoParent creates a feature branch with no defined parent in this repository.
func (r *Runner) CreateFeatureBranchNoParent(name string) error {
	res, err := r.Run("git", "branch", name, "main")
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (r *Runner) CreateFile(name, content string) error {
	filePath := filepath.Join(r.WorkingDir(), name)
	folderPath := filepath.Dir(filePath)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create folder %q: %v", folderPath, err)
	}
	err = ioutil.WriteFile(filePath, []byte(content), 0744)
	if err != nil {
		return fmt.Errorf("cannot create file %q: %w", name, err)
	}
	return nil
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (r *Runner) CreatePerennialBranches(names ...string) error {
	for _, name := range names {
		err := r.CreateBranch(name, "main")
		if err != nil {
			return fmt.Errorf("cannot create perennial branch %q in repo %q: %w", name, r.WorkingDir(), err)
		}
	}
	r.AddToPerennialBranches(names...)
	return nil
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (r *Runner) CreateRemoteBranch(localSha, branchName string) error {
	outcome, err := r.Run("git", "push", "origin", localSha+":refs/heads/"+branchName)
	if err != nil {
		return fmt.Errorf("cannot create remote branch for local SHA %q: %w\n%s", localSha, err, outcome.Output())
	}
	return nil
}

// CreateStandaloneTag creates a tag not on a branch
func (r *Runner) CreateStandaloneTag(name string) error {
	return r.RunMany([][]string{
		{"git", "checkout", "-b", "temp"},
		{"touch", "a.txt"},
		{"git", "add", "-A"},
		{"git", "commit", "-m", "temp"},
		{"git", "tag", "-a", name, "-m", name},
		{"git", "checkout", "-"},
		{"git", "branch", "-D", "temp"},
	})
}

// CreateTag creates a tag with the given name
func (r *Runner) CreateTag(name string) error {
	_, err := r.Run("git", "tag", "-a", name, "-m", name)
	return err
}

// CreateTrackingBranch creates a remote tracking branch for the given local branch.
func (r *Runner) CreateTrackingBranch(branch string) error {
	outcome, err := r.Run("git", "push", "-u", "origin", branch)
	if err != nil {
		return fmt.Errorf("cannot create tracking branch for %q: %w\n%s", branch, err, outcome.Output())
	}
	return nil
}

// CurrentBranch provides the currently checked out branch for this repo.
func (r *Runner) CurrentBranch() (result string, err error) {
	outcome, err := r.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return result, fmt.Errorf("cannot determine the current branch: %w\n%s", err, outcome.Output())
	}
	return outcome.OutputSanitized(), nil
}

// CurrentSha provides the SHA of the currently checked out branch/commit.
func (r *Runner) CurrentSha() (string, error) {
	return r.BranchSha("HEAD")
}

// DeleteLastCommit resets HEAD to the previous commit.
func (r *Runner) DeleteLastCommit() error {
	out, err := r.Run("git", "reset", "--hard", "HEAD~1")
	if err != nil {
		return fmt.Errorf("cannot delete last commit: %w\n%s", err, out.Output())
	}
	return nil
}

// DeleteLocalBranch removes the local branch with the given name.
func (r *Runner) DeleteLocalBranch(name string, force bool) error {
	args := []string{"branch", "-d", name}
	if force {
		args[1] = "-D"
	}
	out, err := r.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot delete local branch %q: %w\n%s", name, err, out.Output())
	}
	return nil
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (r *Runner) DeleteMainBranchConfiguration() error {
	res, err := r.Run("git", "config", "--unset", "git-town.main-branch-name")
	if err != nil {
		return fmt.Errorf("cannot delete main branch configuration: %w\n%s", err, res.Output())
	}
	return nil
}

// DeleteRemoteBranch removes the remote branch of the given local branch.
func (r *Runner) DeleteRemoteBranch(name string) error {
	out, err := r.Run("git", "push", "origin", ":"+name)
	if err != nil {
		return fmt.Errorf("cannot delete tracking branch for %q: %w\n%s", name, err, out.Output())
	}
	return nil
}

// DiscardOpenChanges deletes all uncommitted changes.
func (r *Runner) DiscardOpenChanges() error {
	out, err := r.Run("git", "reset", "--hard")
	if err != nil {
		return fmt.Errorf("cannot discard open changes: %w\n%s", err, out.Output())
	}
	return nil
}

// ExpectedPreviouslyCheckedOutBranch returns what is the expected previously checked out branch
// given the inputs
func (r *Runner) ExpectedPreviouslyCheckedOutBranch(initialPreviouslyCheckedOutBranch, initialBranch string) (string, error) {
	hasInitialPreviouslyCheckedOutBranch, err := r.HasLocalBranch(initialPreviouslyCheckedOutBranch)
	if err != nil {
		return "", err
	}
	if hasInitialPreviouslyCheckedOutBranch {
		currentBranch, err := r.CurrentBranch()
		if err != nil {
			return "", err
		}
		hasInitialBranch, err := r.HasLocalBranch(initialBranch)
		if err != nil {
			return "", err
		}
		if currentBranch == initialBranch || !hasInitialBranch {
			return initialPreviouslyCheckedOutBranch, nil
		}
		return initialBranch, nil
	}
	return Config().GetMainBranch(), nil
}

// Fetch retrieves the updates from the remote repo.
func (r *Runner) Fetch() error {
	_, err := r.Run("git", "fetch")
	if err != nil {
		return fmt.Errorf("cannot fetch: %w", err)
	}
	return nil
}

// FetchUpstream fetches updates from the upstream remote.
func (r *Runner) FetchUpstream(branch string) error {
	out, err := r.Run("git", "fetch", "upstream", branch)
	if err != nil {
		return fmt.Errorf("cannot fetch from upstream: %w\n%s", err, out.Output())
	}
	return nil
}

// FileContent provides the current content of a file.
func (r *Runner) FileContent(filename string) (result string, err error) {
	outcome, err := r.Run("cat", filename)
	if err != nil {
		return result, err
	}
	return outcome.Output(), nil
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (r *Runner) FileContentInCommit(sha string, filename string) (result string, err error) {
	outcome, err := r.Run("git", "show", sha+":"+filename)
	if err != nil {
		return result, fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	return outcome.OutputSanitized(), nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (r *Runner) FilesInCommit(sha string) (result []string, err error) {
	outcome, err := r.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return result, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(outcome.OutputSanitized(), "\n"), nil
}

// FilesInBranch provides the list of the files present in the given branch.
func (r *Runner) FilesInBranch(branch string) (result []string, err error) {
	outcome, err := r.Run("git", "ls-tree", "-r", "--name-only", branch)
	if err != nil {
		return result, fmt.Errorf("cannot determine files in branch %q in repo %q: %w", branch, r.WorkingDir(), err)
	}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
		file := strings.TrimSpace(line)
		if file != "" {
			result = append(result, file)
		}
	}
	return result, err
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their remote
func (r *Runner) HasBranchesOutOfSync() (bool, error) {
	res, err := r.Run("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	if err != nil {
		return false, fmt.Errorf("cannot determine if branches are out of sync in %q: %w %q", r.WorkingDir(), err, res.Output())
	}
	return strings.Contains(res.Output(), "["), nil
}

// HasFile indicates whether this repository contains a file with the given name and content.
func (r *Runner) HasFile(name, content string) (result bool, err error) {
	rawContent, err := ioutil.ReadFile(filepath.Join(r.WorkingDir(), name))
	if err != nil {
		return result, fmt.Errorf("repo doesn't have file %q: %w", name, err)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return result, fmt.Errorf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return true, nil
}

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func (r *Runner) HasGitTownConfigNow() (result bool, err error) {
	outcome, err := r.Run("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		exitError := err.(*exec.ExitError)
		if exitError.ExitCode() == 1 {
			return false, nil
		}
	}
	return outcome.OutputSanitized() != "", err
}

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (r *Runner) HasLocalBranch(name string) (bool, error) {
	branches, err := r.LocalBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local branch %q exists: %w", name, err)
	}
	return util.DoesStringArrayContain(branches, name), nil
}

// HasLocalOrRemoteBranch indicates whether this repo has a local or remote branch with the given name.
func (r *Runner) HasLocalOrRemoteBranch(name string) (bool, error) {
	branches, err := r.LocalAndRemoteBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local or remote branch %q exists: %w", name, err)
	}
	return util.DoesStringArrayContain(branches, name), nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (r *Runner) HasMergeInProgress() (result bool, err error) {
	_, err = os.Stat(filepath.Join(r.WorkingDir(), ".git", "MERGE_HEAD"))
	return err == nil, nil
}

// HasOpenChanges indicates whether this repo has open changes.
func (r *Runner) HasOpenChanges() (bool, error) {
	outcome, err := r.Run("git", "status", "--porcelain")
	if err != nil {
		return false, fmt.Errorf("cannot determine open changes: %w", err)
	}
	return outcome.OutputSanitized() != "", nil
}

// HasRebaseInProgress indicates whether this Git repository currently has a rebase in progress.
func (r *Runner) HasRebaseInProgress() (bool, error) {
	res, err := r.Run("git", "status")
	if err != nil {
		return false, fmt.Errorf("cannot determine rebase in %q progress: %w", r.WorkingDir(), err)
	}
	output := res.OutputSanitized()
	if strings.Contains(output, "You are currently rebasing") {
		return true, nil
	}
	if strings.Contains(output, "rebase in progress") {
		return true, nil
	}
	return false, nil
}

// HasRemote indicates whether this repo has a remote with the given name.
func (r *Runner) HasRemote(name string) (result bool, err error) {
	remotes, err := r.Remotes()
	if err != nil {
		return false, fmt.Errorf("cannot determine if remote %q exists: %w", name, err)
	}
	return util.DoesStringArrayContain(remotes, name), nil
}

// HasShippableChanges indicates whether the given branch has changes
// not currently in the main branch.
func (r *Runner) HasShippableChanges(branch string) (bool, error) {
	out, err := r.Run("git", "diff", r.GetMainBranch()+".."+branch)
	if err != nil {
		return false, fmt.Errorf("cannot determine whether branch %q has shippable changes: %w\n%s", branch, err, out.Output())
	}
	return out.OutputSanitized() != "", nil
}

// HasTrackingBranch indicates whether the local branch with the given name has a remote tracking branch.
func (r *Runner) HasTrackingBranch(name string) (result bool, err error) {
	trackingBranchName := "origin/" + name
	remoteBranches, err := r.RemoteBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine if tracking branch %q exists: %w", name, err)
	}
	for _, line := range remoteBranches {
		if strings.TrimSpace(line) == trackingBranchName {
			return true, nil
		}
	}
	return false, nil
}

// LastActiveDir provides the directory that was last used in this repo.
func (r *Runner) LastActiveDir() (string, error) {
	res, err := r.Run("git", "rev-parse", "--show-toplevel")
	return filepath.FromSlash(res.OutputSanitized()), err
}

// LastCommitMessage returns the commit message for the last commit
func (r *Runner) LastCommitMessage() (string, error) {
	out, err := r.Run("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf("cannot determine last commit message: %w\n%s", err, out.Output())
	}
	return out.OutputSanitized(), nil
}

// LocalBranches provides the names of all local branches in this repo.
func (r *Runner) LocalBranches() (result []string, err error) {
	outcome, err := r.Run("git", "branch")
	if err != nil {
		return result, fmt.Errorf("cannot determine the local branches")
	}
	lines := outcome.OutputLines()
	for l := range lines {
		result = append(result, strings.TrimSpace(strings.Trim(lines[l], "* ")))
	}
	return MainFirst(sort.StringSlice(result)), nil
}

// LocalAndRemoteBranches provides the names of all local branches in this repo.
func (r *Runner) LocalAndRemoteBranches() ([]string, error) {
	outcome, err := r.Run("git", "branch", "-a")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine the local branches")
	}
	lines := outcome.OutputLines()
	branchNames := make(map[string]struct{})
	for l := range lines {
		if !strings.Contains(lines[l], " -> ") {
			branchNames[strings.TrimSpace(strings.Replace(strings.Replace(lines[l], "* ", "", 1), "remotes/origin/", "", 1))] = struct{}{}
		}
	}
	result := make([]string, len(branchNames))
	i := 0
	for branchName := range branchNames {
		result[i] = branchName
		i++
	}
	sort.Strings(result)
	return MainFirst(result), nil
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (r *Runner) MergeBranchNoEdit(branch string) error {
	_, err := r.Run("git", "merge", "--no-edit", branch)
	return err
}

// PreviouslyCheckedOutBranch provides the name of the branch that was previously checked out in this repo.
func (r *Runner) PreviouslyCheckedOutBranch() (name string, err error) {
	outcome, err := r.Run("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return "", fmt.Errorf("cannot determine the previously checked out branch: %w", err)
	}
	return outcome.OutputSanitized(), nil
}

// Pull fetches updates from the origin remote and updates the currently checked out branch.
func (r *Runner) Pull() error {
	outcome, err := r.Run("git", "pull")
	if err != nil {
		return fmt.Errorf("cannot pull updates: %w\n%s", err, outcome.Output())
	}
	return nil
}

// PushBranch pushes the branch with the given name to the remote.
func (r *Runner) PushBranch() error {
	outcome, err := r.Run("git", "push")
	if err != nil {
		return fmt.Errorf("cannot push branch in repo %q to origin: %w\n%v", r.WorkingDir(), err, outcome)
	}
	return nil
}

// PushBranchForce pushes the branch with the given name to the remote.
func (r *Runner) PushBranchForce(name string) error {
	outcome, err := r.Run("git", "push", "-f", "origin", name)
	if err != nil {
		return fmt.Errorf("cannot force-push branch %q in repo %q to origin: %w\n%v", name, r.WorkingDir(), err, outcome)
	}
	return nil
}

// PushBranchSetUpstream pushes the branch with the given name to the remote.
func (r *Runner) PushBranchSetUpstream(name string) error {
	outcome, err := r.Run("git", "push", "-u", "origin", name)
	if err != nil {
		return fmt.Errorf("cannot push branch %q in repo %q to origin: %w\n%v", name, r.WorkingDir(), err, outcome)
	}
	return nil
}

// PushTags pushes new the Git tags to origin.
func (r *Runner) PushTags() error {
	outcome, err := r.Run("git", "push", "--tags")
	if err != nil {
		return fmt.Errorf("cannot push branch in repo %q: %w\n%v", r.WorkingDir(), err, outcome)
	}
	return nil
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (r *Runner) Rebase(target string) error {
	outcome, err := r.Run("git", "rebase", target)
	if err != nil {
		return fmt.Errorf("cannot rebase against branch %q: %w\n%v", target, err, outcome)
	}
	return nil
}

// RemoteBranches provides the names of the remote branches in this repo.
func (r *Runner) RemoteBranches() ([]string, error) {
	if r.remoteBranchCache.Initialized() {
		return r.remoteBranchCache.Get(), nil
	}
	outcome, err := r.Run("git", "branch", "-r")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine remote branches")
	}
	lines := outcome.OutputLines()
	result := make([]string, 0, len(lines)-1)
	for l := range lines {
		if !strings.Contains(lines[l], " -> ") {
			result = append(result, strings.TrimSpace(lines[l]))
		}
	}
	r.remoteBranchCache.Set(result)
	remoteBranchesInitialized = true
	return result, nil
}

// Remotes provides the names of all Git remotes in this repository.
func (r *Runner) Remotes() (names []string, err error) {
	out, err := r.Run("git", "remote")
	if err != nil {
		return names, err
	}
	if out.OutputSanitized() == "" {
		return []string{}, nil
	}
	return out.OutputLines(), nil
}

// RemoveBranch deletes the branch with the given name from this repo.
func (r *Runner) RemoveBranch(name string) error {
	res, err := r.Run("git", "branch", "-D", name)
	if err != nil {
		return fmt.Errorf("cannot delete branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// RemoveRemote deletes the Git remote with the given name.
func (r *Runner) RemoveRemote(name string) error {
	_, err := r.Run("git", "remote", "rm", name)
	return err
}

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func (r *Runner) RemoveUnnecessaryFiles() error {
	fullPath := filepath.Join(r.WorkingDir(), ".git", "hooks")
	err := os.RemoveAll(fullPath)
	if err != nil {
		return fmt.Errorf("cannot remove unnecessary files in %q: %w", fullPath, err)
	}
	_ = os.Remove(filepath.Join(r.WorkingDir(), ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(r.WorkingDir(), ".git", "description"))
	return nil
}

// ShouldPushBranch returns whether the local branch with the given name
// contains commits that have not been pushed to the remote.
func (r *Runner) ShouldPushBranch(branch string) (bool, error) {
	trackingBranch := r.TrackingBranchName(branch)
	out, err := r.Run("git", "rev-list", "--left-right", branch+"..."+trackingBranch)
	if err != nil {
		return false, fmt.Errorf("cannot list diff of %q and %q: %w\n%s", branch, trackingBranch, err, out.Output())
	}
	return out.OutputSanitized() != "", nil
}

// SquashMerge squash-merges the given branch into the current branch
func (r *Runner) SquashMerge(branch string) error {
	out, err := r.Run("git", "merge", "--squash", branch)
	if err != nil {
		return fmt.Errorf("cannot squash-merge branch %q: %w\n%s", branch, err, out.Output())
	}
	return nil
}

// Stash adds the current files to the Git stash.
func (r *Runner) Stash() error {
	err := r.RunMany([][]string{
		{"git", "add", "."},
		{"git", "stash"},
	})
	if err != nil {
		return fmt.Errorf("cannot stash: %w", err)
	}
	return nil
}

// StashSize provides the number of stashes in this repository.
func (r *Runner) StashSize() (result int, err error) {
	res, err := r.Run("git", "stash", "list")
	if err != nil {
		return result, fmt.Errorf("command %q failed: %w", res.FullCmd(), err)
	}
	if res.OutputSanitized() == "" {
		return 0, nil
	}
	return len(res.OutputLines()), nil
}

// Tags provides a list of the tags in this repository
func (r *Runner) Tags() (result []string, err error) {
	res, err := r.Run("git", "tag")
	if err != nil {
		return result, fmt.Errorf("cannot determine tags in repo %q: %w", r.WorkingDir(), err)
	}
	for _, line := range strings.Split(res.OutputSanitized(), "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result, err
}

// TrackingBranchName provides the name of the remote branch tracking the given local branch.
func (r *Runner) TrackingBranchName(branch string) string {
	return "origin/" + branch
}

// UncommittedFiles provides the names of the files not committed into Git.
func (r *Runner) UncommittedFiles() (result []string, err error) {
	res, err := r.Run("git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return result, fmt.Errorf("cannot determine uncommitted files in %q: %w", r.WorkingDir(), err)
	}
	lines := res.OutputLines()
	for l := range lines {
		if lines[l] == "" {
			continue
		}
		parts := strings.Split(lines[l], " ")
		result = append(result, parts[1])
	}
	return result, nil
}

// ShaForCommit provides the SHA for the commit with the given name.
func (r *Runner) ShaForCommit(name string) (result string, err error) {
	var args []string
	if name == "Initial commit" {
		args = []string{"reflog", "--grep=" + name, "--format=%H", "--max-count=1"}
	} else {
		args = []string{"reflog", "--grep-reflog=commit: " + name, "--format=%H"}
	}
	res, err := r.Run("git", args...)
	if err != nil {
		return result, fmt.Errorf("cannot determine SHA of commit %q: %w\n%s", name, err, res.Output())
	}
	if res.OutputSanitized() == "" {
		return result, fmt.Errorf("cannot find the SHA of commit %q", name)
	}
	return res.OutputSanitized(), nil
}

// StageFiles adds the file with the given name to the Git index.
func (r *Runner) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := r.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot stage files %s: %w", strings.Join(names, ", "), err)
	}
	return nil
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (r *Runner) StartCommit() error {
	out, err := r.Run("git", "commit")
	if err != nil {
		return fmt.Errorf("cannot start commit: %w\n%s", err, out.Output())
	}
	return nil
}
