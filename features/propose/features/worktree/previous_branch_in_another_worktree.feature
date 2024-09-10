@skipWindows
Feature: prepend with the previous branch checked out in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exist
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                                            |
      | current | frontend | git fetch --prune --tags                                           |
      | <none>  | frontend | Looking for proposal online ... ok                                 |
      | current | frontend | git checkout main                                                  |
      | main    | frontend | git rebase origin/main                                             |
      |         | frontend | git checkout current                                               |
      | current | frontend | git merge --no-edit --ff main                                      |
      |         | frontend | git push -u origin current                                         |
      | <none>  | frontend | open https://github.com/git-town/git-town/compare/current?expand=1 |
    And the current branch is still "current"
    And the previous Git branch is now "current"
