Feature: dry-run merging branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | beta-file | beta content |
    And the current branch is "beta"
    When I run "git-town merge --dry-run"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                  |
      | beta   | git fetch --prune --tags                 |
      |        | git checkout alpha                       |
      | alpha  | git reset --hard {{ sha 'beta commit' }} |
      |        | git push origin :beta                    |
      |        | git branch -D beta                       |
    And the initial branches exist now
    And the initial commits exist now
  #
  # Cannot test undo because dry-run now doesn't create a runstate.
