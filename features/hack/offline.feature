Feature: git town hack: offline mode

  When having no internet connection
  I want that new branches are created without attempting network accesses
  So that I don't see unnecessary errors.

  Background:
    Given Git Town is in offline mode
    And the following commits exist in my repo
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
    And I am on the "main" branch
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
      | BRANCH  | LOCATION      | MESSAGE     |
      | main    | local, remote | main commit |
