Feature: inside a committed subfolder that exists only on the current feature branch

  Background:
    Given my repo has a feature branch "existing-feature"
    And my repo contains the commits
      | BRANCH           | LOCATION      | MESSAGE       | FILE NAME        |
      | existing-feature | local, remote | folder commit | new_folder/file1 |
    And I am on the "existing-feature" branch
    When I run "git-town hack new-feature" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                     |
      | existing-feature | git fetch --prune --tags    |
      |                  | git checkout main           |
      | main             | git rebase origin/main      |
      |                  | git branch new-feature main |
      |                  | git checkout new-feature    |
    And I am now on the "new-feature" branch
    And my repo is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | existing-feature | main   |
      | new-feature      | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH      | COMMAND                       |
      | new-feature | git checkout main             |
      | main        | git branch -D new-feature     |
      |             | git checkout existing-feature |
    And I am now on the "existing-feature" branch
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy
