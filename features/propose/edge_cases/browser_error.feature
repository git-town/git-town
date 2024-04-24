@skipWindows
Feature: print the URL when the browser crashes

  Background:
    Given the current branch is a feature branch "feature"
    And the origin is "git@github.com:git-town/git-town"
    And tool "open" is broken
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      |         | git checkout main                                                  |
      | main    | git rebase origin/main                                             |
      |         | git checkout feature                                               |
      | feature | git merge --no-edit --ff origin/feature                            |
      |         | git merge --no-edit --ff main                                      |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/compare/feature?expand=1
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
