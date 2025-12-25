@messyoutput
Feature: rename a branch that has an overridden branch type

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE         | PARENT | LOCATIONS     |
      | old  | contribution |        | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And Git setting "git-town-branch.old.branchtype" is "feature"
    And the current branch is "old"
    When I run "git-town rename new" and enter into the dialog:
      | DIALOG                  | KEYS  |
      | parent branch for "old" | enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | old    | git fetch --prune --tags  |
      |        | git branch --move old new |
      |        | git checkout new          |
      | new    | git push -u origin new    |
      |        | git push origin :old      |
    And Git setting "git-town-branch.new.branchtype" is now "feature"
    And Git setting "git-town-branch.old.branchtype" now doesn't exist
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git push -u origin old                |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
      |        | git push origin :new                  |
    And Git setting "git-town-branch.old.branchtype" is now "feature"
    And Git setting "git-town-branch.new.branchtype" now doesn't exist
    And the initial branches and lineage exist now
