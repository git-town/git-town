Feature: git town-hack: push branch to remote upon creation

  Background:
    Given the new-branch-push-flag configuration is true
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE       |
      | main   | remote   | remote commit |
    And I am on the "main" branch
    When I run "git-town hack feature"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | main    | git fetch --prune --tags   |
      |         | git rebase origin/main     |
      |         | git branch feature main    |
      |         | git checkout feature       |
      | feature | git push -u origin feature |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE       |
      | main    | local, remote | remote commit |
      | feature | local, remote | remote commit |
