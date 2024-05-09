Feature: sync a branch when the previous branch is active in another worktree

  Background:
    Given the feature branches "current" and "previous"
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | current | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout current                    |
      | current | git merge --no-edit --ff origin/current |
      |         | git merge --no-edit --ff main           |
    And the current branch is still "current"
    And no commits exist now

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
