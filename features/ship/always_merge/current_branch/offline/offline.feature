Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And offline mode is enabled
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                             |
      | feature | git checkout main                   |
      | main    | git merge --no-ff --edit -- feature |
      |         | git branch -D feature               |
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                |
      | main    | local    | feature commit         |
      |         |          | Merge branch 'feature' |
      | feature | origin   | feature commit         |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And the initial commits exist now
    And the initial branches and lineage exist now
