Feature: a branch with a new remote commit is still rebased onto its tracking branch when pushing is disabled

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | branch | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | branch | local    | local commit  |
      |        | origin   | remote commit |
    And Git setting "git-town.push-branches" is "false"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "branch"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                             |
      | branch | git fetch --prune --tags                            |
      |        | git -c rebase.updateRefs=false rebase origin/branch |
    And no rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | branch | git reset --hard {{ sha 'local commit' }} |
    And the initial commits exist now
