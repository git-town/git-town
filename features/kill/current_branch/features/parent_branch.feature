Feature: delete a branch within a branch chain

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
      | gamma | feature | beta   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "beta" and the previous branch is "alpha"
    And an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                          |
      | beta   | git fetch --prune --tags                         |
      |        | git push origin :beta                            |
      |        | git add -A                                       |
      |        | git commit -m "Committing WIP for git town undo" |
      |        | git checkout alpha                               |
      | alpha  | git branch -D beta                               |
    And it prints:
      """
      branch "gamma" is now a child of "alpha"
      """
    And the current branch is now "alpha"
    And no uncommitted files exist
    And the branches are now
      | REPOSITORY    | BRANCHES           |
      | local, origin | main, alpha, gamma |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | gamma  | local, origin | gamma commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | gamma  | alpha  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                      |
      | alpha  | git push origin {{ sha 'beta commit' }}:refs/heads/beta      |
      |        | git branch beta {{ sha 'Committing WIP for git town undo' }} |
      |        | git checkout beta                                            |
      | beta   | git reset --soft HEAD~1                                      |
    And the current branch is now "beta"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
