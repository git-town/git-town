Feature: sync a branch whose tracking branch was shipped in offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature   | feature | main   | local, origin |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local, origin | feature-1 commit | feature-1-file | feature 1 content |
      | feature-2 | local, origin | feature-2 commit | feature-2-file | feature 2 content |
    And offline mode is enabled
    And origin ships the "feature-1" branch using the "squash-merge" ship-strategy
    And the current branch is "feature-1"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs no commands
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And the initial branches and lineage exist now
