Feature: Print URL when opening browser fails

  When using Git Town on a machine with a broken open browser command
  I want that it prints the URL of the website
  So that I can copy-and-paste it into my browser.

  Background:
    Given my repository has a feature branch named "feature"
    And my repo's origin is "git@github.com:git-town/git-town"
    And my computer has a broken "open" tool installed
    And I am on the "feature" branch
    When I run "git-town new-pull-request"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      |         | git checkout main                                                  |
      | main    | git rebase origin/main                                             |
      |         | git checkout feature                                               |
      | feature | git merge --no-edit origin/feature                                 |
      |         | git merge --no-edit main                                           |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/compare/feature?expand=1
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git checkout main    |
      | main    | git checkout feature |
