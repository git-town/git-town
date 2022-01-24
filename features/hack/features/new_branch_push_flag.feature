Feature: auto-push the new branch

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
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git push origin :feature |
      |         | git checkout main        |
      | main    | git branch -d feature    |
    And I am now on the "main" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE       |
      | main    | local, remote | remote commit |
    And Git Town now has no branch hierarchy information
