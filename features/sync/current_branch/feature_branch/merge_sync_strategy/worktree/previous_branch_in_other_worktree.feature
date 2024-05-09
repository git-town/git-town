Feature: sync a branch when the previous branch is active in another worktree

  Background:
    Given the feature branches "current" and "previous"
    And the commits
      | BRANCH   | LOCATION | MESSAGE                |
      | current  | local    | local current commit   |
      |          | origin   | origin current commit  |
      | previous | local    | local previous commit  |
      |          | origin   | origin previous commit |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                |
      | child  | git fetch --prune --tags               |
      |        | git checkout main                      |
      | main   | git rebase origin/main                 |
      |        | git push                               |
      |        | git checkout child                     |
      | child  | git merge --no-edit --ff origin/child  |
      |        | git merge --no-edit --ff origin/parent |
      |        | git push                               |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE                                                 |
      | main   | local, origin, worktree | origin main commit                                      |
      |        |                         | local main commit                                       |
      | child  | local, origin           | local child commit                                      |
      |        |                         | origin child commit                                     |
      |        |                         | Merge remote-tracking branch 'origin/child' into child  |
      |        |                         | origin parent commit                                    |
      |        |                         | Merge remote-tracking branch 'origin/parent' into child |
      | parent | origin                  | origin parent commit                                    |
      |        | worktree                | local parent commit                                     |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                            |
      | child  | git reset --hard {{ sha 'local child commit' }}                                    |
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
