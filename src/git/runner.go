package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/git-town/git-town/v7/src/stringslice"
)

// Runner executes Git commands.
type Runner struct {
	run.Shell                         // for running console commands
	Config             config.GitTown // caches Git configuration settings
	CurrentBranchCache *cache.String  // caches the currently checked out Git branch
	DryRun             *DryRun        // tracks dry-run information
	IsRepoCache        *cache.Bool    // caches whether the current directory is a Git repo
	RemoteBranchCache  *cache.Strings // caches the remote branches of this Git repo
	RemotesCache       *cache.Strings // caches Git remotes
	RootDirCache       *cache.String  // caches the base of the Git directory
}

// AbortMerge cancels a currently ongoing Git merge operation.
func (r *Runner) AbortMerge() error {
	_, err := r.Run("git", "merge", "--abort")
	if err != nil {
		return fmt.Errorf("cannot abort current merge: %w", err)
	}
	return nil
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (r *Runner) AbortRebase() error {
	_, err := r.Run("git", "rebase", "--abort")
	if err != nil {
		return fmt.Errorf("cannot abort current merge: %w", err)
	}
	return nil
}

// AddRemote adds a Git remote with the given name and URL to this repository.
func (r *Runner) AddRemote(name, url string) error {
	_, err := r.Run("git", "remote", "add", name, url)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w", name, url, err)
	}
	r.RemotesCache.Invalidate()
	return nil
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (r *Runner) AddSubmodule(url string) error {
	_, err := r.Run("git", "submodule", "add", url)
	if err != nil {
		return fmt.Errorf("cannot add submodule %q: %w", url, err)
	}
	return r.Commit("added submodule", "")
}

// Author provides the locally Git configured user.
func (r *Runner) Author() (string, error) {
	out, err := r.Run("git", "config", "user.name")
	if err != nil {
		return "", err
	}
	name := out.OutputSanitized()
	out, err = r.Run("git", "config", "user.email")
	if err != nil {
		return "", err
	}
	email := out.OutputSanitized()
	return name + " <" + email + ">", nil
}

// BranchHasUnmergedCommits indicates whether the branch with the given name
// contains commits that are not merged into the main branch.
func (r *Runner) BranchHasUnmergedCommits(branch string) (bool, error) {
	out, err := r.Run("git", "log", r.Config.MainBranch()+".."+branch)
	if err != nil {
		return false, fmt.Errorf("cannot determine if branch %q has unmerged commits: %w", branch, err)
	}
	return out.OutputSanitized() != "", nil
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (r *Runner) CheckoutBranch(name string) error {
	_, err := r.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q in repo %q: %w", name, r.WorkingDir(), err)
	}
	if name != "-" {
		r.CurrentBranchCache.Set(name)
	} else {
		r.CurrentBranchCache.Invalidate()
	}
	return nil
}

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided.
func (r *Runner) CommentOutSquashCommitMessage(prefix string) error {
	squashMessageFile := ".git/SQUASH_MSG"
	contentBytes, err := os.ReadFile(squashMessageFile)
	if err != nil {
		return fmt.Errorf("cannot read squash message file %q: %w", squashMessageFile, err)
	}
	content := string(contentBytes)
	if prefix != "" {
		content = prefix + "\n" + content
	}
	content = regexp.MustCompile("(?m)^").ReplaceAllString(content, "# ")
	return os.WriteFile(squashMessageFile, []byte(content), 0o600)
}

// CommitNoEdit commits all staged files with the default commit message.
func (r *Runner) CommitNoEdit() error {
	_, err := r.Run("git", "commit", "--no-edit")
	if err != nil {
		return fmt.Errorf("cannot commit files: %w", err)
	}
	return nil
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (r *Runner) Commits(fields []string) ([]Commit, error) {
	branches, err := r.LocalBranchesMainFirst()
	if err != nil {
		return []Commit{}, fmt.Errorf("cannot determine the Git branches: %w", err)
	}
	result := []Commit{}
	for _, branch := range branches {
		commits, err := r.CommitsInBranch(branch, fields)
		if err != nil {
			return []Commit{}, err
		}
		result = append(result, commits...)
	}
	return result, nil
}

// CommitsInBranch provides all commits in the given Git branch.
func (r *Runner) CommitsInBranch(branch string, fields []string) ([]Commit, error) {
	outcome, err := r.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return []Commit{}, fmt.Errorf("cannot get commits in branch %q: %w", branch, err)
	}
	result := []Commit{}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
		parts := strings.Split(line, "|")
		commit := Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if stringslice.Contains(fields, "FILE NAME") {
			filenames, err := r.FilesInCommit(commit.SHA)
			if err != nil {
				return []Commit{}, fmt.Errorf("cannot determine file name for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileName = strings.Join(filenames, ", ")
		}
		if stringslice.Contains(fields, "FILE CONTENT") {
			filecontent, err := r.FileContentInCommit(commit.SHA, commit.FileName)
			if err != nil {
				return []Commit{}, fmt.Errorf("cannot determine file content for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result, nil
}

// CommitStagedChanges commits the currently staged changes.
func (r *Runner) CommitStagedChanges(message string) error {
	var err error
	if message != "" {
		_, err = r.Run("git", "commit", "-m", message)
	} else {
		_, err = r.Run("git", "commit", "--no-edit")
	}
	if err != nil {
		return fmt.Errorf("cannot commit staged changes: %w", err)
	}
	return nil
}

// Commit performs a commit of the staged changes with an optional custom message and author.
func (r *Runner) Commit(message, author string) error {
	gitArgs := []string{"commit"}
	if message != "" {
		gitArgs = append(gitArgs, "-m", message)
	}
	if author != "" {
		gitArgs = append(gitArgs, "--author", author)
	}
	_, err := r.Run("git", gitArgs...)
	return err
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (r *Runner) ConnectTrackingBranch(name string) error {
	_, err := r.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	if err != nil {
		return fmt.Errorf("cannot connect tracking branch for %q: %w", name, err)
	}
	return nil
}

// ContinueRebase continues the currently ongoing rebase.
func (r *Runner) ContinueRebase() error {
	_, err := r.Run("git", "rebase", "--continue")
	if err != nil {
		return fmt.Errorf("cannot continue rebase: %w", err)
	}
	return nil
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (r *Runner) CreateBranch(name, parent string) error {
	_, err := r.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create branch %q: %w", name, err)
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
	_ = r.Config.SetParent(name, parent)
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
	_, err = r.Run("git", "add", commit.FileName)
	if err != nil {
		return fmt.Errorf("cannot add file to commit: %w", err)
	}
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	_, err = r.Run("git", commands...)
	if err != nil {
		return fmt.Errorf("cannot commit: %w", err)
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
	_, err := r.Run("git", "branch", name, "main")
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w", name, err)
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (r *Runner) CreateFile(name, content string) error {
	filePath := filepath.Join(r.WorkingDir(), name)
	folderPath := filepath.Dir(filePath)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create folder %q: %w", folderPath, err)
	}
	err = os.WriteFile(filePath, []byte(content), 0o500)
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
	return r.Config.AddToPerennialBranches(names...)
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (r *Runner) CreateRemoteBranch(localSha, branch string, noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, config.OriginRemote, localSha+":refs/heads/"+branch)
	_, err := r.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot create remote branch for local SHA %q: %w", localSha, err)
	}
	return nil
}

// CreateStandaloneTag creates a tag not on a branch.
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

// CreateTag creates a tag with the given name.
func (r *Runner) CreateTag(name string) error {
	_, err := r.Run("git", "tag", "-a", name, "-m", name)
	return err
}

// CurrentBranch provides the currently checked out branch for this repo.
func (r *Runner) CurrentBranch() (string, error) {
	if r.DryRun.IsActive() {
		return r.DryRun.CurrentBranch(), nil
	}
	if r.CurrentBranchCache.Initialized() {
		return r.CurrentBranchCache.Value(), nil
	}
	rebasing, err := r.HasRebaseInProgress()
	if err != nil {
		return "", fmt.Errorf("cannot determine current branch: %w", err)
	}
	if rebasing {
		currentBranch, err := r.currentBranchDuringRebase()
		if err != nil {
			return "", err
		}
		r.CurrentBranchCache.Set(currentBranch)
		return currentBranch, nil
	}
	outcome, err := r.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("cannot determine the current branch: %w", err)
	}
	r.CurrentBranchCache.Set(outcome.OutputSanitized())
	return r.CurrentBranchCache.Value(), nil
}

func (r *Runner) currentBranchDuringRebase() (string, error) {
	rootDir, err := r.RootDirectory()
	if err != nil {
		return "", err
	}
	rawContent, err := os.ReadFile(fmt.Sprintf("%s/.git/rebase-apply/head-name", rootDir))
	if err != nil {
		// Git 2.26 introduces a new rebase backend, see https://github.com/git/git/blob/master/Documentation/RelNotes/2.26.0.txt
		rawContent, err = os.ReadFile(fmt.Sprintf("%s/.git/rebase-merge/head-name", rootDir))
		if err != nil {
			return "", err
		}
	}
	content := strings.TrimSpace(string(rawContent))
	return strings.ReplaceAll(content, "refs/heads/", ""), nil
}

// CurrentSha provides the SHA of the currently checked out branch/commit.
func (r *Runner) CurrentSha() (string, error) {
	return r.ShaForBranch("HEAD")
}

// DeleteLastCommit resets HEAD to the previous commit.
func (r *Runner) DeleteLastCommit() error {
	_, err := r.Run("git", "reset", "--hard", "HEAD~1")
	if err != nil {
		return fmt.Errorf("cannot delete last commit: %w", err)
	}
	return nil
}

// DeleteLocalBranch removes the local branch with the given name.
func (r *Runner) DeleteLocalBranch(name string, force bool) error {
	args := []string{"branch", "-d", name}
	if force {
		args[1] = "-D"
	}
	_, err := r.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot delete local branch %q: %w", name, err)
	}
	return nil
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (r *Runner) DeleteMainBranchConfiguration() error {
	_, err := r.Run("git", "config", "--unset", config.MainBranchKey)
	if err != nil {
		return fmt.Errorf("cannot delete main branch configuration: %w", err)
	}
	return nil
}

// DeleteRemoteBranch removes the remote branch of the given local branch.
func (r *Runner) DeleteRemoteBranch(name string) error {
	_, err := r.Run("git", "push", config.OriginRemote, ":"+name)
	if err != nil {
		return fmt.Errorf("cannot delete tracking branch for %q: %w", name, err)
	}
	return nil
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (r *Runner) DiffParent(branch, parentBranch string) error {
	_, err := r.Run("git", "diff", parentBranch+".."+branch)
	if err != nil {
		return fmt.Errorf("cannot diff branch %q with its parent branch %q: %w", branch, parentBranch, err)
	}
	return nil
}

// DiscardOpenChanges deletes all uncommitted changes.
func (r *Runner) DiscardOpenChanges() error {
	_, err := r.Run("git", "reset", "--hard")
	if err != nil {
		return fmt.Errorf("cannot discard open changes: %w", err)
	}
	return nil
}

// ExpectedPreviouslyCheckedOutBranch returns what is the expected previously checked out branch
// given the inputs.
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
	return r.Config.MainBranch(), nil
}

// Fetch retrieves the updates from the origin repo.
func (r *Runner) Fetch() error {
	_, err := r.Run("git", "fetch", "--prune", "--tags")
	if err != nil {
		return fmt.Errorf("cannot fetch: %w", err)
	}
	return nil
}

// FetchUpstream fetches updates from the upstream remote.
func (r *Runner) FetchUpstream(branch string) error {
	_, err := r.Run("git", "fetch", "upstream", branch)
	if err != nil {
		return fmt.Errorf("cannot fetch from upstream: %w", err)
	}
	return nil
}

// FileContent provides the current content of a file.
func (r *Runner) FileContent(filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(r.WorkingDir(), filename))
	return string(content), err
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (r *Runner) FileContentInCommit(sha string, filename string) (string, error) {
	outcome, err := r.Run("git", "show", sha+":"+filename)
	if err != nil {
		return "", fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	result := outcome.OutputSanitized()
	if strings.HasPrefix(result, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		result = ""
	}
	return result, nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (r *Runner) FilesInCommit(sha string) ([]string, error) {
	outcome, err := r.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return []string{}, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(outcome.OutputSanitized(), "\n"), nil
}

// FilesInBranch provides the list of the files present in the given branch.
func (r *Runner) FilesInBranch(branch string) ([]string, error) {
	outcome, err := r.Run("git", "ls-tree", "-r", "--name-only", branch)
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine files in branch %q in repo %q: %w", branch, r.WorkingDir(), err)
	}
	result := []string{}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
		file := strings.TrimSpace(line)
		if file != "" {
			result = append(result, file)
		}
	}
	return result, err
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their tracking branch.
func (r *Runner) HasBranchesOutOfSync() (bool, error) {
	res, err := r.Run("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	if err != nil {
		return false, fmt.Errorf("cannot determine if branches are out of sync in %q: %w %q", r.WorkingDir(), err, res.Output())
	}
	return strings.Contains(res.Output(), "["), nil
}

// HasConflicts returns whether the local repository currently has unresolved merge conflicts.
func (r *Runner) HasConflicts() (bool, error) {
	res, err := r.Run("git", "status")
	if err != nil {
		return false, fmt.Errorf("cannot determine conflicts: %w", err)
	}
	return res.OutputContainsText("Unmerged paths"), nil
}

// HasFile indicates whether this repository contains a file with the given name and content.
func (r *Runner) HasFile(name, content string) (bool, error) {
	rawContent, err := os.ReadFile(filepath.Join(r.WorkingDir(), name))
	if err != nil {
		return false, fmt.Errorf("repo doesn't have file %q: %w", name, err)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return false, fmt.Errorf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return true, nil
}

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func (r *Runner) HasGitTownConfigNow() (bool, error) {
	outcome, err := r.Run("git", "config", "--local", "--get-regex", "git-town")
	if outcome.ExitCode() == 1 {
		return false, nil
	}
	return outcome.OutputSanitized() != "", err
}

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (r *Runner) HasLocalBranch(name string) (bool, error) {
	branches, err := r.LocalBranchesMainFirst()
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local branch %q exists: %w", name, err)
	}
	return stringslice.Contains(branches, name), nil
}

// HasLocalOrRemoteBranch indicates whether this repo or origin have a branch with the given name.
func (r *Runner) HasLocalOrOriginBranch(name string) (bool, error) {
	branches, err := r.LocalAndOriginBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local or remote branch %q exists: %w", name, err)
	}
	return stringslice.Contains(branches, name), nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (r *Runner) HasMergeInProgress() (bool, error) {
	_, err := os.Stat(filepath.Join(r.WorkingDir(), ".git", "MERGE_HEAD"))
	return err == nil, nil
}

// HasOpenChanges indicates whether this repo has open changes.
func (r *Runner) HasOpenChanges() (bool, error) {
	outcome, err := r.Run("git", "status", "--porcelain", "--ignore-submodules")
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

// HasOrigin indicates whether this repo has an origin remote.
func (r *Runner) HasOrigin() (bool, error) {
	return r.HasRemote(config.OriginRemote)
}

// HasRemote indicates whether this repo has a remote with the given name.
func (r *Runner) HasRemote(name string) (bool, error) {
	remotes, err := r.Remotes()
	if err != nil {
		return false, fmt.Errorf("cannot determine if remote %q exists: %w", name, err)
	}
	return stringslice.Contains(remotes, name), nil
}

// HasShippableChanges indicates whether the given branch has changes
// not currently in the main branch.
func (r *Runner) HasShippableChanges(branch string) (bool, error) {
	out, err := r.Run("git", "diff", r.Config.MainBranch()+".."+branch)
	if err != nil {
		return false, fmt.Errorf("cannot determine whether branch %q has shippable changes: %w", branch, err)
	}
	return out.OutputSanitized() != "", nil
}

// HasTrackingBranch indicates whether the local branch with the given name has a remote tracking branch.
func (r *Runner) HasTrackingBranch(name string) (bool, error) {
	trackingBranch := "origin/" + name
	remoteBranches, err := r.RemoteBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine if tracking branch %q exists: %w", name, err)
	}
	for _, line := range remoteBranches {
		if strings.TrimSpace(line) == trackingBranch {
			return true, nil
		}
	}
	return false, nil
}

// IsBranchInSync returns whether the branch with the given name is in sync with its tracking branch.
func (r *Runner) IsBranchInSync(branch string) (bool, error) {
	hasTrackingBranch, err := r.HasTrackingBranch(branch)
	if err != nil {
		return false, err
	}
	if hasTrackingBranch {
		localSha, err := r.ShaForBranch(branch)
		if err != nil {
			return false, err
		}
		remoteSha, err := r.ShaForBranch(r.TrackingBranch(branch))
		return localSha == remoteSha, err
	}
	return true, nil
}

// IsRepository returns whether or not the current directory is in a repository.
func (r *Runner) IsRepository() bool {
	if !r.IsRepoCache.Initialized() {
		_, err := run.Exec("git", "rev-parse")
		r.IsRepoCache.Set(err == nil)
	}
	return r.IsRepoCache.Value()
}

// LastCommitMessage provides the commit message for the last commit.
func (r *Runner) LastCommitMessage() (string, error) {
	out, err := r.Run("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf("cannot determine last commit message: %w", err)
	}
	return out.OutputSanitized(), nil
}

// LocalAndOriginBranches provides the names of all local branches in this repo.
func (r *Runner) LocalAndOriginBranches() ([]string, error) {
	outcome, err := r.Run("git", "branch", "-a")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine the local branches")
	}
	lines := outcome.OutputLines()
	branch := make(map[string]struct{})
	for _, line := range lines {
		if !strings.Contains(line, " -> ") {
			branch[strings.TrimSpace(strings.Replace(strings.Replace(line, "* ", "", 1), "remotes/origin/", "", 1))] = struct{}{}
		}
	}
	result := make([]string, len(branch))
	i := 0
	for branch := range branch {
		result[i] = branch
		i++
	}
	sort.Strings(result)
	mainBranch := r.Config.MainBranchOr("main")
	return stringslice.Hoist(result, mainBranch), nil
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
func (r *Runner) LocalBranches() ([]string, error) {
	res, err := r.Run("git", "branch")
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range res.OutputLines() {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return result, nil
}

// LocalBranchesMainFirst provides the names of all local branches in this repo.
func (r *Runner) LocalBranchesMainFirst() ([]string, error) {
	branches, err := r.LocalBranches()
	if err != nil {
		return []string{}, err
	}
	mainBranch := r.Config.MainBranchOr("main")
	return stringslice.Hoist(sort.StringSlice(branches), mainBranch), nil
}

// LocalBranchesWithDeletedTrackingBranches provides the names of all branches
// whose remote tracking branches have been deleted.
func (r *Runner) LocalBranchesWithDeletedTrackingBranches() ([]string, error) {
	res, err := r.Run("git", "branch", "-vv")
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range res.OutputLines() {
		line = strings.Trim(line, "* ")
		parts := strings.SplitN(line, " ", 2)
		branch := parts[0]
		deleteTrackingBranchStatus := fmt.Sprintf("[%s: gone]", r.TrackingBranch(branch))
		if strings.Contains(parts[1], deleteTrackingBranchStatus) {
			result = append(result, branch)
		}
	}
	return result, nil
}

// LocalBranchesWithoutMain provides the names of all branches in the local repository,
// ordered alphabetically without the main branch.
func (r *Runner) LocalBranchesWithoutMain() ([]string, error) {
	mainBranch := r.Config.MainBranch()
	branches, err := r.LocalBranches()
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, branch := range branches {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return result, nil
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (r *Runner) MergeBranchNoEdit(branch string) error {
	_, err := r.Run("git", "merge", "--no-edit", branch)
	return err
}

// PopStash restores stashed-away changes into the workspace.
func (r *Runner) PopStash() error {
	_, err := r.Run("git", "stash", "pop")
	if err != nil {
		return fmt.Errorf("cannot pop the stash: %w", err)
	}
	return nil
}

// PreviouslyCheckedOutBranch provides the name of the branch that was previously checked out in this repo.
func (r *Runner) PreviouslyCheckedOutBranch() (string, error) {
	outcome, err := r.Run("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return "", fmt.Errorf("cannot determine the previously checked out branch: %w", err)
	}
	return outcome.OutputSanitized(), nil
}

// Pull fetches updates from origin and updates the currently checked out branch.
func (r *Runner) Pull() error {
	_, err := r.Run("git", "pull")
	if err != nil {
		return fmt.Errorf("cannot pull updates: %w", err)
	}
	return nil
}

type PushArgs struct {
	Branch         string
	Force          bool `exhaustruct:"optional"`
	ForceWithLease bool `exhaustruct:"optional"`
	NoPushHook     bool `exhaustruct:"optional"`
	Remote         string
}

// PushBranch pushes the branch with the given name to origin.
func (r *Runner) PushBranch(options ...PushArgs) error {
	var option PushArgs
	if len(options) > 0 {
		option = options[0]
	} else {
		option = PushArgs{} //nolint:exhaustruct  // intentional zero-value object
	}
	args := []string{"push"}
	provideBranch := false
	if option.Force {
		args = append(args, "-f")
	}
	if option.NoPushHook {
		args = append(args, "--no-verify")
	}
	if option.ForceWithLease {
		args = append(args, "--force-with-lease")
	}
	if option.Remote != "" {
		args = append(args, "-u", option.Remote)
		provideBranch = true
	}
	if option.Branch != "" && provideBranch {
		args = append(args, option.Branch)
	}
	_, err := r.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot push branch in repo %q to origin: %w", r.WorkingDir(), err)
	}
	return nil
}

// PushTags pushes new the Git tags to origin.
func (r *Runner) PushTags() error {
	_, err := r.Run("git", "push", "--tags")
	if err != nil {
		return fmt.Errorf("cannot push branch in repo %q: %w", r.WorkingDir(), err)
	}
	return nil
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (r *Runner) Rebase(target string) error {
	_, err := r.Run("git", "rebase", target)
	if err != nil {
		return fmt.Errorf("cannot rebase against branch %q: %w", target, err)
	}
	return nil
}

// RemoteBranches provides the names of the remote branches in this repo.
func (r *Runner) RemoteBranches() ([]string, error) {
	if !r.RemoteBranchCache.Initialized() {
		outcome, err := r.Run("git", "branch", "-r")
		if err != nil {
			return []string{}, fmt.Errorf("cannot determine remote branches: %w", err)
		}
		lines := outcome.OutputLines()
		branches := make([]string, 0, len(lines)-1)
		for _, line := range lines {
			if !strings.Contains(line, " -> ") {
				branches = append(branches, strings.TrimSpace(line))
			}
		}
		r.RemoteBranchCache.Set(branches)
	}
	return r.RemoteBranchCache.Value(), nil
}

// Remotes provides the names of all Git remotes in this repository.
func (r *Runner) Remotes() ([]string, error) {
	if !r.RemotesCache.Initialized() {
		out, err := r.Run("git", "remote")
		if err != nil {
			return []string{}, fmt.Errorf("cannot determine remotes: %w", err)
		}
		if out.OutputSanitized() == "" {
			r.RemotesCache.Set([]string{})
		} else {
			r.RemotesCache.Set(out.OutputLines())
		}
	}
	return r.RemotesCache.Value(), nil
}

// RemoveBranch deletes the branch with the given name from this repo.
func (r *Runner) RemoveBranch(name string) error {
	_, err := r.Run("git", "branch", "-D", name)
	if err != nil {
		return fmt.Errorf("cannot delete branch %q: %w", name, err)
	}
	return nil
}

// RemoveRemote deletes the Git remote with the given name.
func (r *Runner) RemoveRemote(name string) error {
	r.RemotesCache.Invalidate()
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

// ResetToSha undoes all commits on the current branch all the way until the given SHA.
func (r *Runner) ResetToSha(sha string, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha)
	_, err := r.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot reset to SHA %q: %w", sha, err)
	}
	return nil
}

// RevertCommit reverts the commit with the given SHA.
func (r *Runner) RevertCommit(sha string) error {
	_, err := r.Run("git", "revert", sha)
	if err != nil {
		return fmt.Errorf("cannot revert commit %q: %w", sha, err)
	}
	return nil
}

// RootDirectory provides the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func (r *Runner) RootDirectory() (string, error) {
	if !r.RootDirCache.Initialized() {
		res, err := r.Run("git", "rev-parse", "--show-toplevel")
		if err != nil {
			return "", fmt.Errorf("cannot determine root directory: %w", err)
		}
		r.RootDirCache.Set(filepath.FromSlash(res.OutputSanitized()))
	}
	return r.RootDirCache.Value(), nil
}

// ShaForBranch provides the SHA for the local branch with the given name.
func (r *Runner) ShaForBranch(name string) (string, error) {
	outcome, err := r.Run("git", "rev-parse", name)
	if err != nil {
		return "", fmt.Errorf("cannot determine SHA of local branch %q: %w", name, err)
	}
	return outcome.OutputSanitized(), nil
}

// ShaForCommit provides the SHA for the commit with the given name.
func (r *Runner) ShaForCommit(name string) (string, error) {
	res, err := r.Run("git", "log", "--reflog", "--format=%H", "--grep=^"+name+"$")
	if err != nil {
		return "", fmt.Errorf("cannot determine the SHA of commit %q: %w", name, err)
	}
	result := res.OutputSanitized()
	if result == "" {
		return "", fmt.Errorf("cannot find the SHA of commit %q", name)
	}
	result = strings.Split(result, "\n")[0]
	return result, nil
}

// ShouldPushBranch returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (r *Runner) ShouldPushBranch(branch string) (bool, error) {
	trackingBranch := r.TrackingBranch(branch)
	out, err := r.Run("git", "rev-list", "--left-right", branch+"..."+trackingBranch)
	if err != nil {
		return false, fmt.Errorf("cannot list diff of %q and %q: %w", branch, trackingBranch, err)
	}
	return out.OutputSanitized() != "", nil
}

// SquashMerge squash-merges the given branch into the current branch.
func (r *Runner) SquashMerge(branch string) error {
	_, err := r.Run("git", "merge", "--squash", branch)
	if err != nil {
		return fmt.Errorf("cannot squash-merge branch %q: %w", branch, err)
	}
	return nil
}

// Stash adds the current files to the Git stash.
func (r *Runner) Stash() error {
	err := r.RunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
	if err != nil {
		return fmt.Errorf("cannot stash: %w", err)
	}
	return nil
}

// StashSize provides the number of stashes in this repository.
func (r *Runner) StashSize() (int, error) {
	res, err := r.Run("git", "stash", "list")
	if err != nil {
		return 0, fmt.Errorf("command %q failed: %w", res.FullCmd(), err)
	}
	if res.OutputSanitized() == "" {
		return 0, nil
	}
	return len(res.OutputLines()), nil
}

// Tags provides a list of the tags in this repository.
func (r *Runner) Tags() ([]string, error) {
	res, err := r.Run("git", "tag")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine tags in repo %q: %w", r.WorkingDir(), err)
	}
	result := []string{}
	for _, line := range strings.Split(res.OutputSanitized(), "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result, err
}

// TrackingBranch provides the name of the remote branch tracking the local branch with the given name.
func (r *Runner) TrackingBranch(branch string) string {
	return "origin/" + branch
}

// UncommittedFiles provides the names of the files not committed into Git.
func (r *Runner) UncommittedFiles() ([]string, error) {
	res, err := r.Run("git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine uncommitted files in %q: %w", r.WorkingDir(), err)
	}
	lines := res.OutputLines()
	result := []string{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		result = append(result, parts[1])
	}
	return result, nil
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
	_, err := r.Run("git", "commit")
	if err != nil {
		return fmt.Errorf("cannot start commit: %w", err)
	}
	return nil
}

// Version indicates whether the needed Git version is installed.
//
//nolint:nonamedreturns  // multiple int return values justify using names for return values
func (r *Runner) Version() (major int, minor int, err error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	res, err := r.Run("git", "version")
	if err != nil {
		return 0, 0, fmt.Errorf("cannot determine Git version: %w", err)
	}
	matches := versionRegexp.FindStringSubmatch(res.OutputSanitized())
	if matches == nil {
		return 0, 0, fmt.Errorf("'git version' returned unexpected output: %q.\nPlease open an issue and supply the output of running 'git version'", res.Output())
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert major version %q to int: %w", matches[1], err)
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert minor version %q to int: %w", matches[2], err)
	}
	return majorVersion, minorVersion, nil
}
