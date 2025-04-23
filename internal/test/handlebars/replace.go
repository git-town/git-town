package handlebars

import (
	"fmt"
	"log"
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

var templateOnce sync.Once

func Replace(text string, localRepo, remoteRepo, worktreeRepo Runner, initialDevSHAs map[string]gitdomain.SHA, initialOriginSHAsOpt, initialWorktreeSHAsOpt Option[map[string]gitdomain.SHA]) string {
	var templateRE *regexp.Regexp
	templateOnce.Do(func() { templateRE = regexp.MustCompile(`\{\{.*?\}\}`) })
	for strings.Contains(text, "{{") {
		templateOnce.Do(func() { templateRE = regexp.MustCompile(`\{\{.*?\}\}`) })
		match := templateRE.FindString(text)
		switch {
		case strings.HasPrefix(match, "{{ sha "):
			commitName := match[8 : len(match)-4]
			shas := localRepo.SHAsForCommit(commitName)
			if len(shas) == 0 {
				panic(fmt.Sprintf("test workspace has no commit %q", commitName))
			}
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-origin "):
			commitName := match[18 : len(match)-4]
			shas := remoteRepo.SHAsForCommit(commitName)
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-before-run "):
			commitName := match[19 : len(match)-4]
			sha, found := initialDevSHAs[commitName]
			if !found {
				fmt.Printf("I cannot find the initial dev commit %q.\n", commitName)
				fmt.Printf("I have records about %d commits:\n", len(initialDevSHAs))
				for key := range maps.Keys(initialDevSHAs) {
					fmt.Println("  -", key)
				}
				panic("see error above")
			}
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-origin-before-run "):
			initialOriginSHAs, has := initialOriginSHAsOpt.Get()
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
			shas := worktreeRepo.SHAsForCommit(commitName)
			sha := shas.First()
			text = strings.Replace(text, match, sha.String(), 1)
		case strings.HasPrefix(match, "{{ sha-in-worktree-before-run "):
			commitName := match[31 : len(match)-4]
			initialWorktreeSHAs, has := initialWorktreeSHAsOpt.Get()
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
			log.Fatalf("DataTable.Expand: unknown template expression %q", text)
		}
	}
	return text
}
