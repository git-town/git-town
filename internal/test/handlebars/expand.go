package handlebars

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type GitCommands interface {
	SHAsForCommit(gitdomain.CommitMessage) gitdomain.SHAs
}

var (
	templateOnce sync.Once
	templateRE   *regexp.Regexp
)

func Expand(text string, args ExpandArgs) string {
	expandOnce.Do(func() {
		expandRegex = regexp.MustCompile(`\{\{.*?\}\}`)
	})
	templateOnce.Do(func() { templateRE = expandRegex })
	for strings.Contains(text, "{{") {
		match := templateRE.FindString(text)
		switch {
		case strings.HasPrefix(match, "{{ sha "):
			commitMessage := gitdomain.CommitMessage(match[8 : len(match)-4])
			shas := args.LocalRepo.SHAsForCommit(commitMessage)
			if len(shas) == 0 {
				panic(fmt.Sprintf("test workspace has no commit %q", commitMessage))
			}
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-short "):
			commitMessage := gitdomain.CommitMessage(match[14 : len(match)-4])
			shas := args.LocalRepo.SHAsForCommit(commitMessage)
			if len(shas) == 0 {
				panic(fmt.Sprintf("test workspace has no commit %q", commitMessage))
			}
			sha := shas.First().Truncate(7)
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-origin "):
			commitMessage := gitdomain.CommitMessage(match[18 : len(match)-4])
			shas := args.RemoteRepo.SHAsForCommit(commitMessage)
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-initial "):
			commitMessage := gitdomain.CommitMessage(match[16 : len(match)-4])
			commit, hasCommit := args.InitialDevCommits.FindByCommitMessage(commitMessage).Get()
			if !hasCommit {
				fmt.Printf("I cannot find the initial dev commit %q.\n", commitMessage)
				fmt.Printf("I have records about %d commits:\n", len(args.InitialDevCommits))
				for _, commit := range args.InitialDevCommits {
					fmt.Printf("  - %q (%s)\n", commit.Message, commit.SHA)
				}
				panic("see error above")
			}
			text = strings.Replace(text, match, commit.SHA.String(), 1)
		case strings.HasPrefix(match, "{{ sha-initial-short "):
			commitMessage := gitdomain.CommitMessage(match[22 : len(match)-4])
			commit, hasCommit := args.InitialDevCommits.FindByCommitMessage(commitMessage).Get()
			if !hasCommit {
				fmt.Printf("I cannot find the initial dev commit %q.\n", commitMessage)
				fmt.Printf("I have records about %d commits:\n", len(args.InitialDevCommits))
				for _, commit := range args.InitialDevCommits {
					fmt.Printf("  - %q (%s)\n", commit.Message, commit.SHA)
				}
				panic("see error above")
			}
			text = strings.Replace(text, match, commit.SHA.Truncate(7).String(), 1)
		case strings.HasPrefix(match, "{{ sha-before-run "):
			commitMessage := gitdomain.CommitMessage(match[19 : len(match)-4])
			commit, found := args.BeforeRunDevSHAs.FindByCommitMessage(commitMessage).Get()
			if !found {
				fmt.Printf("I cannot find the before-run dev commit %q.\n", commitMessage)
				fmt.Printf("I have records about %d commits:\n", len(args.BeforeRunDevSHAs))
				for _, commit := range args.BeforeRunDevSHAs {
					fmt.Printf("  - %q (%s)\n", commit.Message, commit.SHA)
				}
				panic("see error above")
			}
			text = strings.Replace(text, match, commit.SHA.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-origin-initial "):
			initialOriginCommits, has := args.InitialOriginCommits.Get()
			if !has {
				panic("no origin SHAs recorded")
			}
			commitMessage := gitdomain.CommitMessage(match[26 : len(match)-4])
			commit, hasCommit := initialOriginCommits.FindByCommitMessage(commitMessage).Get()
			if !hasCommit {
				fmt.Printf("I cannot find the initial origin commit %q.\n", commitMessage)
				fmt.Printf("I have records about %d commits:\n", len(initialOriginCommits))
				for _, commit := range initialOriginCommits {
					fmt.Printf("  - %q (%s)\n", commit.Message, commit.SHA)
				}
			}
			text = strings.Replace(text, match, commit.SHA.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-origin-before-run "):
			beforeRunOriginSHAs, has := args.BeforeRunOriginSHAsOpt.Get()
			if !has {
				panic("no origin SHAs recorded")
			}
			commitMessage := gitdomain.CommitMessage(match[29 : len(match)-4])
			commit, hasCommit := beforeRunOriginSHAs.FindByCommitMessage(commitMessage).Get()
			if !hasCommit {
				fmt.Printf("I cannot find the initial origin commit %q.\n", commitMessage)
				fmt.Printf("I have records about %d commits:\n", len(beforeRunOriginSHAs))
				for _, commit := range beforeRunOriginSHAs {
					fmt.Printf("  - %q (%s)\n", commit.Message, commit.SHA)
				}
			}
			text = strings.Replace(text, match, commit.SHA.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-worktree "):
			commitMessage := gitdomain.CommitMessage(match[20 : len(match)-4])
			shas := args.WorktreeRepo.SHAsForCommit(commitMessage)
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-worktree-initial "):
			commitMessage := gitdomain.CommitMessage(match[28 : len(match)-4])
			initialWorktreeSHAs, has := args.InitialWorktreeCommits.Get()
			if !has {
				panic("no initial worktree SHAs recorded")
			}
			commit, hasCommit := initialWorktreeSHAs.FindByCommitMessage(commitMessage).Get()
			if !hasCommit {
				fmt.Printf("I cannot find the initial worktree commit %q.\n", commitMessage)
				fmt.Printf("I have records about %d commits:\n", len(initialWorktreeSHAs))
				for _, commit := range initialWorktreeSHAs {
					fmt.Printf("  - %q (%s)\n", commit.Message, commit.SHA)
				}
			}
			text = strings.Replace(text, match, commit.SHA.String(), 1)
		default:
			panic(fmt.Sprintf("DataTable.Expand: unknown template expression %q", match))
		}
	}
	return text
}

var (
	expandOnce  sync.Once
	expandRegex *regexp.Regexp
)

type ExpandArgs struct {
	BeforeRunDevSHAs       gitdomain.Commits
	BeforeRunOriginSHAsOpt Option[gitdomain.Commits]
	InitialDevCommits      gitdomain.Commits
	InitialOriginCommits   Option[gitdomain.Commits]
	InitialWorktreeCommits Option[gitdomain.Commits]
	LocalRepo              GitCommands
	RemoteRepo             GitCommands
	WorktreeRepo           GitCommands
}
