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
    And the current branch is "prototype"
    And Git Town setting "sync-prototype-strategy" is "rebase"
    And Git Town setting "sync-feature-strategy" is "merge"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND         |
      | prototype | git rebase main |
    And these commits exist now
      | BRANCH    | LOCATION | MESSAGE      |
      | main      | local    | main commit  |
      | prototype | local    | main commit  |
      |           |          | local commit |
    And all branches are now synchronized
    And the current branch is still "prototype"
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                                              |
      | prototype | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "prototype"
    And the initial commits exist now
    And the initial branches and lineage exist now
