Feature: auto-push the new branch

  Background:
    Given a Git repo with origin
    And Git setting "git-town.push-new-branches" is "true"
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | main   | origin   | origin commit |
    And the current branch is "main"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | main   | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout -b new         |
      | new    | git push -u origin new      |
      |        | git stash pop               |
    And the current branch is now "new"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | new    | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout main           |
      | main   | git branch -D new           |
      |        | git push origin :new        |
      |        | git stash pop               |
    And the current branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
