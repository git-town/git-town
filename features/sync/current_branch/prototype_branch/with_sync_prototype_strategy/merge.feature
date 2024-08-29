Feature: sync the current prototype branch in a local repo

  Background:
    Given a local Git repo
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | prototype | prototype | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE      | FILE NAME  |
      | main      | local    | main commit  | main_file  |
      | prototype | local    | local commit | local_file |
    And the current branch is "prototype"
    And Git Town setting "sync-prototype-strategy" is "merge"
    And Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | prototype | git merge --no-edit --ff main |
    And these commits exist now
      | BRANCH    | LOCATION | MESSAGE                            |
      | main      | local    | main commit                        |
      | prototype | local    | local commit                       |
      |           |          | main commit                        |
      |           |          | Merge branch 'main' into prototype |
    And all branches are now synchronized
    And the current branch is still "prototype"
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                                              |
      | prototype | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "prototype"
    And the initial commits exist
    And the initial branches and lineage exist
