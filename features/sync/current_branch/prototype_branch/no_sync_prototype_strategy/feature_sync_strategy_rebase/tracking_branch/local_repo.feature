Feature: sync the current prototype branch in a local repo

  Background:
    Given a local Git repo
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | prototype | prototype | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE      | FILE NAME  |
      | main      | local    | main commit  | main_file  |
      | prototype | local    | local commit | local_file |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "prototype"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                      |
      | prototype | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
    And the initial branches and lineage exist now
    And all branches are now synchronized
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | prototype | git reset --hard {{ sha-initial 'local commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
