Feature: rename a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME  | TYPE         | PARENT | LOCATIONS     |
      | other | contribution | main   | local, origin |
    And the current branch is "other"
    And the commits
      | BRANCH | LOCATION      | MESSAGE             |
      | other  | local, origin | experimental commit |
    When I run "git-town rename-branch other new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
      |        | git branch new other     |
      |        | git checkout new         |
      | new    | git push -u origin new   |
      |        | git push origin :other   |
      |        | git branch -D other      |
    And the current branch is now "new"
    And the contribution branches are now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | new    | local, origin | experimental commit |
    And this lineage exists now
      | BRANCH | PARENT |

  Scenario: undo
    Given I ran "git-town rename-branch --force other new"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                          |
      | new    | git branch other {{ sha 'experimental commit' }} |
      |        | git push -u origin other                         |
      |        | git push origin :new                             |
      |        | git checkout other                               |
      | other  | git branch -D new                                |
    And the current branch is now "other"
    And the contribution branches are now "other"
    And the initial commits exist
    And the initial branches and lineage exist
