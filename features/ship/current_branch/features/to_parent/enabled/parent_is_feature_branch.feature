Feature: allowing shiping into a feature branch

  Background:
    Given a local Git repo
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | alpha | feature | main   | local     |
      | beta  | feature | alpha  | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | alpha  | local    | alpha commit |
      | beta   | local    | beta commit  |
    And the current branch is "beta"
    When I run "git-town ship --to-parent -m done"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                      |
      | beta   | git checkout alpha           |
      | alpha  | git merge --squash --ff beta |
      |        | git commit -m done           |
      |        | git branch -D beta           |
    And the current branch is now "alpha"
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, alpha |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | alpha  | local    | alpha commit |
      |        |          | done         |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | alpha  | git reset --hard {{ sha 'alpha commit' }} |
      |        | git branch beta {{ sha 'beta commit' }}   |
      |        | git checkout beta                         |
    And the current branch is now "beta"
    And the initial commits exist
    And the initial branches and lineage exist
