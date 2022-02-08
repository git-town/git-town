Feature: offline mode

  Background:
    Given Git Town is in offline mode
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And my workspace has an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git add -A             |
      |        | git stash              |
      |        | git rebase origin/main |
      |        | git branch new main    |
      |        | git checkout new       |
      | new    | git stash pop          |
    And I am now on the "new" branch
    And my workspace still contains my uncommitted file
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | main commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git add -A        |
      |        | git stash         |
      |        | git checkout main |
      | main   | git branch -d new |
      |        | git stash pop     |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And now the initial commits exist
    And Git Town is now aware of no branch hierarchy
