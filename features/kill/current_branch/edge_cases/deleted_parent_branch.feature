Feature: the parent of the branch to kill was deleted remotely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And origin deletes the "alpha" branch
    And the current branch is "beta" and the previous branch is "alpha"
    When I run "git-town kill"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
      |        | git push origin :beta    |
      |        | git checkout main        |
      | main   | git branch -D alpha      |
    And the current branch is now "main"
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                     |
      | other  | git branch old {{ sha 'Committing WIP for git town undo' }} |
      |        | git checkout old                                            |
      | old    | git reset --soft HEAD~1                                     |
    And the current branch is now "old"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | old    | local         | old commit   |
      | other  | local, origin | other commit |
    And the uncommitted file still exists
    And the initial branches and lineage exist now
