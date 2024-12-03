Feature: the parent of the branch to delete was deleted remotely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And origin deletes the "alpha" branch
    And the current branch is "beta" and the previous branch is "alpha"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
      |        | git push origin :beta    |
      |        | git checkout alpha       |
      | alpha  | git branch -D beta       |
    And the current branch is now "alpha"
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, alpha |
      | origin     | main        |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                    |
      | alpha  | git branch beta {{ sha 'initial commit' }} |
      |        | git push -u origin beta                    |
      |        | git checkout beta                          |
    And the current branch is now "beta"
    And the initial branches and lineage exist now
