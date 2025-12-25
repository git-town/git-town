Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And offline mode is enabled
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git checkout main                                                            |
      | main    | git -c rebase.updateRefs=false rebase origin/main                            |
      |         | git checkout feature                                                         |
      | feature | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                   |
      | feature | git reset --hard {{ sha-initial 'local feature commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
