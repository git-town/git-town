@skipWindows
Feature: prepend with the previous branch checked out in another worktree

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the current branch is "current" and the previous branch is "previous"
    And branch "previous" is active in another worktree
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | TYPE     | COMMAND                                                            |
      | current | frontend | git fetch --prune --tags                                           |
      |         | frontend | git push -u origin current                                         |
      |         | frontend | Finding proposal from current into main ... ok                     |
      |         | frontend | open https://github.com/git-town/git-town/compare/current?expand=1 |
    And the previous Git branch is now "current"
