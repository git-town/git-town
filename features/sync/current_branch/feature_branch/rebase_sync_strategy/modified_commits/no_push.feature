Feature: a rewritten local branch is not rebased onto its old tracking branch when pushing is disabled

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | branch | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | branch | local, origin | old     |
    And Git setting "git-town.push-branches" is "false"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And I amend this commit
      | BRANCH | LOCATION | MESSAGE |
      | branch | local    | new     |
    And the current branch is "branch"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | branch | git fetch --prune --tags |
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
