<<<<<<< HEAD:features/hack/offline.feature
Feature: git town hack: offline mode
=======
Feature: offline mode
>>>>>>> main:features/hack/features/offline.feature

  Background:
    Given Git Town is in offline mode
    And the following commits exist in my repo
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
<<<<<<< HEAD:features/hack/offline.feature
      | BRANCH  | LOCATION      | MESSAGE     |
      | main    | local, remote | main commit |
=======
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
>>>>>>> main:features/hack/features/offline.feature
    And Git Town now has no branch hierarchy information
