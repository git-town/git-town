Feature: dry-run merging branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
      | beta   | local, origin | beta commit  | beta-file  | beta content  |
    And the current branch is "beta"
    When I run "git-town merge --dry-run"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git branch -D alpha                   |
      |        | git push origin :alpha                |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial branches exist now

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs no commands
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial branches exist now
