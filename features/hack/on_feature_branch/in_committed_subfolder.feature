Feature: git town-hack: creating a feature branch from a subfolder not on the main branch

  To ensure robustness
  When creating a feature branch from a subfolder that doesn't exist on the main branch
  I want that Git Town changes to the root directory before changing branches.

  Background:
    Given my repo has a feature branch named "feature1"
    And the following commits exist in my repo
      | BRANCH   | LOCATION      | MESSAGE       | FILE NAME        |
      | feature1 | local, remote | folder commit | new_folder/file1 |
    And I am on the "feature1" branch
    And my workspace has an uncommitted file
    When I run "git-town hack feature2" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                  |
      | feature1    | git fetch --prune --tags |
      |             | git add -A               |
      |             | git stash                |
      |             | git checkout main        |
      | main        | git rebase origin/main   |
      |             | git branch feature2 main |
      |             | git checkout feature2    |
      | feature2    | git stash pop            |
    And I am now on the "feature2" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH   | PARENT |
      | feature1 | main   |
      | feature2 | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                |
      | feature2 | git add -A             |
      |          | git stash              |
      |          | git checkout main      |
      | main     | git branch -d feature2 |
      |          | git checkout feature1  |
      | feature1 | git stash pop          |
    And I am now on the "feature1" branch
    And my repo is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH   | PARENT |
      | feature1 | main   |
