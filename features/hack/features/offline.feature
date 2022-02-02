Feature: offline mode

  Background:
    Given Git Town is in offline mode
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
    And my workspace has an uncommitted file
    When I run "git-town hack feature"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                 |
      | main    | git add -A              |
      |         | git stash               |
      |         | git rebase origin/main  |
      |         | git branch feature main |
      |         | git checkout feature    |
      | feature | git stash pop           |
    And I am now on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE     |
      | main    | local, remote | main commit |
      | feature | local         | main commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | feature | git add -A            |
      |         | git stash             |
      |         | git checkout main     |
      | main    | git branch -d feature |
      |         | git stash pop         |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
    And Git Town now has no branch hierarchy information
