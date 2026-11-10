Feature: a rewritten local branch is not rebased onto its old tracking branch when pushing is disabled

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE |
      | parent | local, origin | old     |
    And Git setting "git-town.push-branches" is "false"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "parent"
    And I ran "git commit --amend -m new"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | parent | git fetch --prune --tags |
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
