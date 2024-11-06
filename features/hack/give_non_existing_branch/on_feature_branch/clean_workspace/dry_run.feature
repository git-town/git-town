Feature: dry-run hacking a new feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the current branch is "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    When I run "git-town hack new --dry-run"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                 |
      | existing | git fetch --prune --tags                |
      |          | git checkout main                       |
      | main     | git rebase origin/main --no-update-refs |
      |          | git checkout -b new                     |
    And the current branch is still "existing"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "existing"
    And the initial commits exist now
    And the initial branches and lineage exist now
