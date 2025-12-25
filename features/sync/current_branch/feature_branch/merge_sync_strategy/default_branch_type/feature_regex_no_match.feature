@smoke
Feature: a unknown branch type is set, the feature-regex does not match

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE   | PARENT | LOCATIONS     |
      | my-branch | (none) |        | local, origin |
    And the commits
      | BRANCH    | LOCATION | MESSAGE                 |
      | main      | local    | local main commit       |
      |           | origin   | origin main commit      |
      | my-branch | local    | local my-branch commit  |
      |           | origin   | origin my-branch commit |
    And local Git setting "git-town.feature-regex" is "other"
    And local Git setting "git-town.unknown-branch-type" is "observed"
    And the current branch is "my-branch"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                |
      | my-branch | git fetch --prune --tags                               |
      |           | git -c rebase.updateRefs=false rebase origin/my-branch |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                 |
      | main      | local         | local main commit       |
      |           | origin        | origin main commit      |
      | my-branch | local, origin | origin my-branch commit |
      |           | local         | local my-branch commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                             |
      | my-branch | git reset --hard {{ sha 'local my-branch commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
