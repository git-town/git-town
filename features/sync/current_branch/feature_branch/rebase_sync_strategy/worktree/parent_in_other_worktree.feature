Feature: sync a branch whose parent is active in another worktree

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And branch "parent" is active in another worktree
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main                          |
      |        | git push                                        |
      |        | git checkout child                              |
      | child  | git rebase origin/parent                        |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/child                         |
      |        | git push --force-with-lease --force-if-includes |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE              |
      | main   | local, origin, worktree | origin main commit   |
      |        |                         | local main commit    |
      | child  | local, origin           | origin child commit  |
      |        |                         | origin parent commit |
      |        |                         | local child commit   |
      | parent | origin                  | origin parent commit |
      |        | worktree                | local parent commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                            |
      | child  | git reset --hard {{ sha-before-run 'local child commit' }}                         |
      |        | git push --force-with-lease origin {{ sha-in-origin 'origin child commit' }}:child |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE              |
      | main   | local, origin, worktree | origin main commit   |
      |        |                         | local main commit    |
      | child  | local                   | local child commit   |
      |        | origin                  | origin child commit  |
      | parent | origin                  | origin parent commit |
      |        | worktree                | local parent commit  |
