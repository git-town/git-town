Feature: swapping a branch with its prototype parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE      | PARENT | LOCATIONS     |
      | parent  | prototype | main   | local, origin |
      | current | feature   | parent | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | parent  | local, origin | parent commit  |
      | current | local, origin | current commit |
    And the current branch is "current"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                   |
      | current | git fetch --prune --tags                                  |
      |         | git -c rebase.updateRefs=false rebase --onto main parent  |
      |         | git checkout parent                                       |
      | parent  | git -c rebase.updateRefs=false rebase --onto current main |
      |         | git push --force-with-lease --force-if-includes           |
      |         | git checkout current                                      |
    And this lineage exists now
      """
      main
        current
          parent
      """
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | parent  | local, origin | parent commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | current | git checkout parent                             |
      | parent  | git reset --hard {{ sha 'parent commit' }}      |
      |         | git push --force-with-lease --force-if-includes |
      |         | git checkout current                            |
    And the initial lineage exists now
    And the initial commits exist now
