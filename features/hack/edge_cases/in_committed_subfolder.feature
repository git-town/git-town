Feature: inside a committed subfolder that exists only on the current feature branch

  Background:
    Given my repo has a feature branch "existing"
    And my repo contains the commits
      | BRANCH   | LOCATION      | MESSAGE       | FILE NAME        |
      | existing | local, remote | folder commit | new_folder/file1 |
    And I am on the "existing" branch
    When I run "git-town hack new" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
      |          | git branch new main      |
      |          | git checkout new         |
    And I am now on the "new" branch
    And my repo is left with my initial commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH | COMMAND               |
      | new    | git checkout main     |
      | main   | git branch -d new     |
      |        | git checkout existing |
    And I am now on the "existing" branch
    And my repo is left with my initial commits
    And Git Town is now aware of the initial branch hierarchy
