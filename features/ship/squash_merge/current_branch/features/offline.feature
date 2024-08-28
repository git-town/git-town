Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And offline mode is enabled
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m "feature done"    |
      |         | git branch -D feature           |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE        |
      | main    | local    | feature done   |
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
