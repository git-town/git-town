Feature: on a feature branch

  Background:
    Given my repo contains the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | remote   | main_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git fetch --prune --tags    |
      |             | git add -A                  |
      |             | git stash                   |
      |             | git rebase origin/main      |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE     |
      | main        | local, remote | main_commit |
      | new-feature | local         | main_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH      | PARENT |
      | new-feature | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH      | COMMAND                   |
      | new-feature | git add -A                |
      |             | git stash                 |
      |             | git checkout main         |
      | main        | git branch -D new-feature |
      |             | git stash pop             |
    And I am now on the "main" branch
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main_commit |
    And Git Town now has no branch hierarchy information
