Feature: a branch with a new remote commit is still rebased onto its tracking branch when pushing is disabled

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | parent | local    | local commit  |
      |        | origin   | remote commit |
    And Git setting "git-town.push-branches" is "false"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "parent"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                             |
      | parent | git fetch --prune --tags                            |
      |        | git -c rebase.updateRefs=false rebase origin/parent |
    And no rebase is now in progress
