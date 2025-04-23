package handlebars

import (
	"fmt"
	"maps"
	"regexp"
	"strings"
	"sync"

	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

type Runner interface {
	SHAsForCommit(name string) gitdomain.SHAs
}

var (
	templateOnce sync.Once
	templateRE   *regexp.Regexp
)

func Expand(text string, args ExpandArgs) string {
	templateOnce.Do(func() { templateRE = regexp.MustCompile(`\{\{.*?\}\}`) })
	for strings.Contains(text, "{{") {
		match := templateRE.FindString(text)
		switch {
		case strings.HasPrefix(match, "{{ sha "):
			commitName := match[8 : len(match)-4]
			shas := args.LocalRepo.SHAsForCommit(commitName)
			if len(shas) == 0 {
				panic(fmt.Sprintf("test workspace has no commit %q", commitName))
			}
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-short "):
			commitName := match[14 : len(match)-4]
			shas := args.LocalRepo.SHAsForCommit(commitName)
			if len(shas) == 0 {
				panic(fmt.Sprintf("test workspace has no commit %q", commitName))
			}
			sha := shas.First()
			sha = gitdomain.NewSHA(sha.String()[:7])
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-origin "):
			commitName := match[18 : len(match)-4]
			shas := args.RemoteRepo.SHAsForCommit(commitName)
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-before-run "):
			commitName := match[19 : len(match)-4]
			sha, found := args.InitialDevSHAs[commitName]
			if !found {
				fmt.Printf("I cannot find the initial dev commit %q.\n", commitName)
				fmt.Printf("I have records about %d commits:\n", len(args.InitialDevSHAs))
				for key := range maps.Keys(args.InitialDevSHAs) {
					fmt.Println("  -", key)
				}
				panic("see error above")
			}
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-origin-before-run "):
			initialOriginSHAs, has := args.InitialOriginSHAsOpt.Get()
			if !has {
				panic("no origin SHAs recorded")
			}
			commitName := match[29 : len(match)-4]
			sha, found := initialOriginSHAs[commitName]
			if !found {
				fmt.Printf("I cannot find the initial origin commit %q.\n", commitName)
				fmt.Printf("I have records about %d commits:\n", len(initialOriginSHAs))
				for key := range maps.Keys(initialOriginSHAs) {
					fmt.Println("  -", key)
				}
			}
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-worktree "):
			commitName := match[20 : len(match)-4]
			shas := args.WorktreeRepo.SHAsForCommit(commitName)
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-worktree-before-run "):
			commitName := match[31 : len(match)-4]
			initialWorktreeSHAs, has := args.InitialWorktreeSHAsOpt.Get()
			if !has {
				panic("no initial worktree SHAs recorded")
			}
			sha, found := initialWorktreeSHAs[commitName]
			if !found {
				fmt.Printf("I cannot find the initial worktree commit %q.\n", commitName)
				fmt.Printf("I have records about %d commits:\n", len(initialWorktreeSHAs))
				for key := range maps.Keys(initialWorktreeSHAs) {
					fmt.Println("  -", key)
				}
			}
			text = strings.Replace(text, match, sha.String(), 1)
		default:
			panic(fmt.Sprintf("DataTable.Expand: unknown template expression %q", match))
		}
	}
	return text
}

type ExpandArgs struct {
	InitialDevSHAs         map[string]gitdomain.SHA
	InitialOriginSHAsOpt   Option[map[string]gitdomain.SHA]
	InitialWorktreeSHAsOpt Option[map[string]gitdomain.SHA]
	LocalRepo              Runner
	RemoteRepo             Runner
	WorktreeRepo           Runner
}
