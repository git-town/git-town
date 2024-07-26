Feature: auto-push the new branch to origin

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And Git Town setting "push-new-branches" is "true"
    And the current branch is "main"
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git checkout -b new      |
      | new    | git push -u origin new   |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | main commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | new    | git push origin :new                        |
      |        | git checkout main                           |
      | main   | git reset --hard {{ sha 'initial commit' }} |
      |        | git branch -D new                           |
    And the current branch is now "main"
    And the initial commits exist
    And the initial branches and lineage exist
