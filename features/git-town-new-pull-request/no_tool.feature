Feature: print URL when unable to open browser

  As a developer using Git Town on a remote machine via SSH
  I want that the tool prints the URL it would have opened
  So that I can copy-and-paste the URL to open it on my local computer.

  Background:
    Given my repository has a feature branch named "feature"
    And my repo's origin is "git@github.com:git-town/git-town"
    And my computer has no tool to open browsers installed
    And I am on the "feature" branch
    When I run "git-town new-pull-request"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
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
