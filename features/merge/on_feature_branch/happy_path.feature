Feature: merging a branch in a stack with its parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
    And the current branch is "beta"
    And Git Town setting "sync-feature-strategy" is "merge"
    When I run "git-town merge"

  # @debug
  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | beta   | git merge --no-edit --ff alpha |
      |        | git branch -D alpha            |
    And the current branch is still "beta"
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | feature | local, origin | the feature | file      | content 3    |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the current branch is now "existing"
    And the initial commits exist now
    And the initial lineage exists now
