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
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git branch -D feature       |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE        |
      | main    | local    | feature commit |
      | feature | origin   | feature commit |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And the current branch is now "feature"
    And the initial commits exist
    And the initial branches and lineage exist
