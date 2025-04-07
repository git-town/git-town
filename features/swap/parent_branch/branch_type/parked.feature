Feature: swapping a branch with its parked parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | parent  | parked  | main   | local, origin |
      | current | feature | parent | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | parent  | local, origin | parent commit  |
      | current | local, origin | current commit |
    And the current branch is "current"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | current | git fetch --prune --tags                        |
      |         | git rebase --onto main parent --no-update-refs  |
      |         | git checkout parent                             |
      | parent  | git rebase --onto current main --no-update-refs |
      |         | git push --force-with-lease --force-if-includes |
      |         | git checkout current                            |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | parent  | local, origin | parent commit  |
    And this lineage exists now
      | BRANCH  | PARENT  |
      | current | main    |
      | parent  | current |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | current | git checkout parent                             |
      | parent  | git reset --hard {{ sha 'parent commit' }}      |
      |         | git push --force-with-lease --force-if-includes |
      |         | git checkout current                            |
    And the initial commits exist now
    And the initial lineage exists now
