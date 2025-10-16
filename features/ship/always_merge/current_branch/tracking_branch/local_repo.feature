Feature: ship a feature branch in a local repo using the always-merge strategy

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "feature"
    When I run "git-town ship" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                             |
      | feature | git checkout main                   |
      | main    | git merge --no-ff --edit -- feature |
      |         | git branch -D feature               |
    And no lineage exists now
    And the branches are now
      | REPOSITORY | BRANCHES |
      | local      | main     |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE                |
      | main   | local    | feature commit         |
      |        |          | Merge branch 'feature' |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
