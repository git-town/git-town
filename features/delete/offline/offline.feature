Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And offline mode is enabled
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | other   | local, origin | other commit   |
    And the current branch is "feature"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout other    |
      | other   | git branch -D feature |
    And no uncommitted files exist now
    And the branches are now
      | REPOSITORY | BRANCHES             |
      | local      | main, other          |
      | origin     | main, feature, other |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | origin        | feature commit |
      | other   | local, origin | other commit   |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And the initial commits exist now
    And the initial branches and lineage exist now
