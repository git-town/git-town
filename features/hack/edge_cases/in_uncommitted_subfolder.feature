Feature: inside an uncommitted subfolder on the current feature branch

  Background:
    Given my repo has a feature branch "existing"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
    And I am on the "existing" branch
    And my workspace has an uncommitted file in folder "new_folder"
    When I run "git-town hack new" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git add -A               |
      |          | git stash                |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
      |          | git branch new main      |
      |          | git checkout new         |
      | new      | git stash pop            |
    And I am now on the "new" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
      | new    | local         | main commit |
    And Git Town now knows about this branch hierarchy
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout main     |
      | main     | git branch -d new     |
      |          | git checkout existing |
      | existing | git stash pop         |
    And I am now on the "existing" branch
    And my repo is left with my initial commits
    And Git Town now has the initial branch hierarchy
