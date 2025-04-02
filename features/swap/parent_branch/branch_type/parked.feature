Feature: swapping a branch with its parked parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | parked  | parked  | main   | local, origin |
      | feature | feature | parked | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | parked  | local, origin | parked commit  |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git rebase --onto main parked                   |
      |         | git checkout parked                             |
      | parked  | git rebase --onto feature main                  |
      |         | git push --force-with-lease --force-if-includes |
      |         | git checkout feature                            |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | parked  | local, origin | parked commit  |
    And this lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | parked  | feature |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git checkout parked                             |
      | parked  | git reset --hard {{ sha 'parked commit' }}      |
      |         | git push --force-with-lease --force-if-includes |
      |         | git checkout feature                            |
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial lineage exists now
