package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/stringslice"
	"github.com/git-town/git-town/v7/src/subshell"
)

type BackendRunner interface {
	Run(executable string, args ...string) (*subshell.Output, error)
	RunMany([][]string) error
}

// BackendCommands are Git commands that Git Town executes in its backend,
// i.e. invisible to the user. They determine the state of the repo without changing the repo.
type BackendCommands struct {
	BackendRunner             // executes shell commands in the directory of the Git repo
	Config        *RepoConfig // the known state of the Git repository
}

// Author provides the locally Git configured user.
func (r *BackendCommands) Author() (string, error) {
	out, err := r.Run("git", "config", "user.name")
	if err != nil {
		return "", err
	}
	name := out.Sanitized()
	out, err = r.Run("git", "config", "user.email")
	if err != nil {
		return "", err
	}
	email := out.Sanitized()
	return name + " <" + email + ">", nil
}

// BranchAuthors provides the user accounts that contributed to the given branch.
// Returns lines of "name <email>".
func (r *BackendCommands) BranchAuthors(branch, parent string) ([]string, error) {
	lines, err := r.Run("git", "shortlog", "-s", "-n", "-e", parent+".."+branch)
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range lines.Lines() {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, parts[1])
	}
	return result, nil
}

// BranchHasUnmergedCommits indicates whether the branch with the given name
// contains commits that are not merged into the main branch.
func (r *BackendCommands) BranchHasUnmergedCommits(branch, parent string) (bool, error) {
	out, err := r.Run("git", "log", parent+".."+branch)
	if err != nil {
		return false, fmt.Errorf("cannot determine if branch %q has unmerged commits: %w", branch, err)
	}
	return out.Sanitized() != "", nil
}

// CheckoutBranch checks out the Git branch with the given name.
func (r *BackendCommands) CheckoutBranch(name string) error {
	if !r.Config.DryRun {
		_, err := r.Run("git", "checkout", name)
		if err != nil {
			return fmt.Errorf("cannot check out branch %q: %w", name, err)
		}
	}
	if name != "-" {
		r.Config.CurrentBranchCache.Set(name)
	} else {
		r.Config.CurrentBranchCache.Invalidate()
	}
	return nil
}

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided.
func (r *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {
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

// CreateFeatureBranch creates a feature branch with the given name in this repository.
func (r *BackendCommands) CreateFeatureBranch(name string) error {
	err := r.RunMany([][]string{
		{"git", "branch", name, "main"},
		{"git", "config", "git-town-branch." + name + ".parent", "main"},
	})
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w", name, err)
	}
	return nil
}

// CurrentBranch provides the currently checked out branch.
func (r *BackendCommands) CurrentBranch() (string, error) {
	if r.Config.DryRun {
		return r.Config.CurrentBranchCache.Value(), nil
	}
	if r.Config.CurrentBranchCache.Initialized() {
		return r.Config.CurrentBranchCache.Value(), nil
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
		r.Config.CurrentBranchCache.Set(currentBranch)
		return currentBranch, nil
	}
	output, err := r.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("cannot determine the current branch: %w", err)
	}
	r.Config.CurrentBranchCache.Set(output.Sanitized())
	return r.Config.CurrentBranchCache.Value(), nil
}

func (r *BackendCommands) currentBranchDuringRebase() (string, error) {
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
func (r *BackendCommands) CurrentSha() (string, error) {
	return r.ShaForBranch("HEAD")
}

// ExpectedPreviouslyCheckedOutBranch returns what is the expected previously checked out branch
// given the inputs.
func (r *BackendCommands) ExpectedPreviouslyCheckedOutBranch(initialPreviouslyCheckedOutBranch, initialBranch, mainBranch string) (string, error) {
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
	return mainBranch, nil
}

// HasConflicts returns whether the local repository currently has unresolved merge conflicts.
func (r *BackendCommands) HasConflicts() (bool, error) {
	output, err := r.Run("git", "status")
	if err != nil {
		return false, fmt.Errorf("cannot determine conflicts: %w", err)
	}
	return output.ContainsText("Unmerged paths"), nil
}

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (r *BackendCommands) HasLocalBranch(name string) (bool, error) {
	branches, err := r.LocalBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local branch %q exists: %w", name, err)
	}
	return stringslice.Contains(branches, name), nil
}

// HasLocalOrRemoteBranch indicates whether this repo or origin have a branch with the given name.
func (r *BackendCommands) HasLocalOrOriginBranch(name, mainBranch string) (bool, error) {
	branches, err := r.LocalAndOriginBranches(mainBranch)
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local or remote branch %q exists: %w", name, err)
	}
	return stringslice.Contains(branches, name), nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (r *BackendCommands) HasMergeInProgress() bool {
	_, err := r.Run("git", "rev-parse", "-q", "--verify", "MERGE_HEAD")
	return err == nil
}

// HasOpenChanges indicates whether this repo has open changes.
func (r *BackendCommands) HasOpenChanges() (bool, error) {
	output, err := r.Run("git", "status", "--porcelain", "--ignore-submodules")
	if err != nil {
		return false, fmt.Errorf("cannot determine open changes: %w", err)
	}
	return output.Sanitized() != "", nil
}

// HasRebaseInProgress indicates whether this Git repository currently has a rebase in progress.
func (r *BackendCommands) HasRebaseInProgress() (bool, error) {
	output, err := r.Run("git", "status")
	if err != nil {
		return false, fmt.Errorf("cannot determine rebase in progress: %w", err)
	}
	sanitizedOutput := output.Sanitized()
	if strings.Contains(sanitizedOutput, "You are currently rebasing") {
		return true, nil
	}
	if strings.Contains(sanitizedOutput, "rebase in progress") {
		return true, nil
	}
	return false, nil
}

// HasOrigin indicates whether this repo has an origin remote.
func (r *BackendCommands) HasOrigin() (bool, error) {
	return r.HasRemote(config.OriginRemote)
}

// HasRemote indicates whether this repo has a remote with the given name.
func (r *BackendCommands) HasRemote(name string) (bool, error) {
	remotes, err := r.Remotes()
	if err != nil {
		return false, fmt.Errorf("cannot determine if remote %q exists: %w", name, err)
	}
	return stringslice.Contains(remotes, name), nil
}

// HasShippableChanges indicates whether the given branch has changes
// not currently in the main branch.
func (r *BackendCommands) HasShippableChanges(branch, mainBranch string) (bool, error) {
	out, err := r.Run("git", "diff", mainBranch+".."+branch)
	if err != nil {
		return false, fmt.Errorf("cannot determine whether branch %q has shippable changes: %w", branch, err)
	}
	return out.Sanitized() != "", nil
}

// HasTrackingBranch indicates whether the local branch with the given name has a remote tracking branch.
func (r *BackendCommands) HasTrackingBranch(name string) (bool, error) {
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
func (r *BackendCommands) IsBranchInSync(branch string) (bool, error) {
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
func (r *BackendCommands) IsRepository() bool {
	if !r.Config.IsRepoCache.Initialized() {
		_, err := r.Run("git", "rev-parse")
		r.Config.IsRepoCache.Set(err == nil)
	}
	return r.Config.IsRepoCache.Value()
}

// LastCommitMessage provides the commit message for the last commit.
func (r *BackendCommands) LastCommitMessage() (string, error) {
	out, err := r.Run("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf("cannot determine last commit message: %w", err)
	}
	return out.Sanitized(), nil
}

// LocalAndOriginBranches provides the names of all local branches in this repo.
func (r *BackendCommands) LocalAndOriginBranches(mainBranch string) ([]string, error) {
	output, err := r.Run("git", "branch", "-a")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine the local branches")
	}
	lines := output.Lines()
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
	return stringslice.Hoist(result, mainBranch), nil
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
func (r *BackendCommands) LocalBranches() ([]string, error) {
	output, err := r.Run("git", "branch")
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range output.Lines() {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return result, nil
}

// LocalBranchesMainFirst provides the names of all local branches in this repo.
func (r *BackendCommands) LocalBranchesMainFirst(mainBranch string) ([]string, error) {
	branches, err := r.LocalBranches()
	if err != nil {
		return []string{}, err
	}
	return stringslice.Hoist(sort.StringSlice(branches), mainBranch), nil
}

// LocalBranchesWithDeletedTrackingBranches provides the names of all branches
// whose remote tracking branches have been deleted.
func (r *BackendCommands) LocalBranchesWithDeletedTrackingBranches() ([]string, error) {
	output, err := r.Run("git", "branch", "-vv")
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range output.Lines() {
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
func (r *BackendCommands) LocalBranchesWithoutMain(mainBranch string) ([]string, error) {
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

// PreviouslyCheckedOutBranch provides the name of the branch that was previously checked out in this repo.
func (r *BackendCommands) PreviouslyCheckedOutBranch() (string, error) {
	output, err := r.Run("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return "", fmt.Errorf("cannot determine the previously checked out branch: %w", err)
	}
	return output.Sanitized(), nil
}

// RemoteBranches provides the names of the remote branches in this repo.
func (r *BackendCommands) RemoteBranches() ([]string, error) {
	if !r.Config.RemoteBranchCache.Initialized() {
		output, err := r.Run("git", "branch", "-r")
		if err != nil {
			return []string{}, fmt.Errorf("cannot determine remote branches: %w", err)
		}
		lines := output.Lines()
		branches := make([]string, 0, len(lines)-1)
		for _, line := range lines {
			if !strings.Contains(line, " -> ") {
				branches = append(branches, strings.TrimSpace(line))
			}
		}
		r.Config.RemoteBranchCache.Set(branches)
	}
	return r.Config.RemoteBranchCache.Value(), nil
}

// Remotes provides the names of all Git remotes in this repository.
func (r *BackendCommands) Remotes() ([]string, error) {
	if !r.Config.RemotesCache.Initialized() {
		out, err := r.Run("git", "remote")
		if err != nil {
			return []string{}, fmt.Errorf("cannot determine remotes: %w", err)
		}
		if out.Sanitized() == "" {
			r.Config.RemotesCache.Set([]string{})
		} else {
			r.Config.RemotesCache.Set(out.Lines())
		}
	}
	return r.Config.RemotesCache.Value(), nil
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (r *BackendCommands) RemoveOutdatedConfiguration() error {
	branches, err := r.LocalAndOriginBranches(r.Config.MainBranch())
	if err != nil {
		return err
	}
	for child, parent := range r.Config.ParentBranchMap() {
		hasChildBranch := stringslice.Contains(branches, child)
		hasParentBranch := stringslice.Contains(branches, parent)
		if !hasChildBranch || !hasParentBranch {
			err = r.Config.RemoveParent(child)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RootDirectory provides the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func (r *BackendCommands) RootDirectory() (string, error) {
	if !r.Config.RootDirCache.Initialized() {
		output, err := r.Run("git", "rev-parse", "--show-toplevel")
		if err != nil {
			return "", fmt.Errorf("cannot determine root directory: %w", err)
		}
		r.Config.RootDirCache.Set(filepath.FromSlash(output.Sanitized()))
	}
	return r.Config.RootDirCache.Value(), nil
}

// ShaForBranch provides the SHA for the local branch with the given name.
func (r *BackendCommands) ShaForBranch(name string) (string, error) {
	output, err := r.Run("git", "rev-parse", name)
	if err != nil {
		return "", fmt.Errorf("cannot determine SHA of local branch %q: %w", name, err)
	}
	return output.Sanitized(), nil
}

// ShouldPushBranch returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (r *BackendCommands) ShouldPushBranch(branch string) (bool, error) {
	trackingBranch := r.TrackingBranch(branch)
	out, err := r.Run("git", "rev-list", "--left-right", branch+"..."+trackingBranch)
	if err != nil {
		return false, fmt.Errorf("cannot list diff of %q and %q: %w", branch, trackingBranch, err)
	}
	return out.Sanitized() != "", nil
}

// TrackingBranch provides the name of the remote branch tracking the local branch with the given name.
func (r *BackendCommands) TrackingBranch(branch string) string {
	return "origin/" + branch
}

// Version indicates whether the needed Git version is installed.
//
//nolint:nonamedreturns  // multiple int return values justify using names for return values
func (r *BackendCommands) Version() (major int, minor int, err error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	output, err := r.Run("git", "version")
	if err != nil {
		return 0, 0, fmt.Errorf("cannot determine Git version: %w", err)
	}
	matches := versionRegexp.FindStringSubmatch(output.Sanitized())
	if matches == nil {
		return 0, 0, fmt.Errorf("'git version' returned unexpected output: %q.\nPlease open an issue and supply the output of running 'git version'", output.Sanitized())
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
