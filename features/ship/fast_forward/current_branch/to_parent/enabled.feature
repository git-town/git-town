Feature: allowing shipping into a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | alpha  | local, origin | alpha 1 | alpha_1   | alpha 1      |
      |        |               | alpha 2 | alpha_2   | alpha 2      |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta 1  | beta_1    | beta 1       |
      |        |               | beta 2  | beta_2    | beta 2       |
    And the current branch is "beta"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship --to-parent"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
      |        | git checkout alpha       |
      | alpha  | git merge --ff-only beta |
      |        | git push                 |
      |        | git push origin :beta    |
      |        | git branch -D beta       |
    And the current branch is now "alpha"
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, alpha |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE |
      | alpha  | local, origin | alpha 1 |
      |        |               | alpha 2 |
      |        |               | beta 1  |
      |        |               | beta 2  |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git reset --hard {{ sha 'alpha 2' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch beta {{ sha 'beta 2' }}              |
      |        | git push -u origin beta                         |
      |        | git checkout beta                               |
    And the current branch is now "beta"
    And the initial commits exist
    And the initial branches and lineage exist
