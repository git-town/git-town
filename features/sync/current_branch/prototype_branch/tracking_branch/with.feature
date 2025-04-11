Feature: sync the current prototype branch that has a tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION | MESSAGE       | FILE NAME   |
      | main      | local    | main commit   | main_file   |
      | prototype | local    | local commit  | local_file  |
      |           | origin   | origin commit | origin_file |
    And the current branch is "prototype"
    And Git setting "git-town.sync-prototype-strategy" is "rebase"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | prototype | git fetch --prune --tags                        |
      |           | git checkout main                               |
      | main      | git rebase origin/main --no-update-refs         |
      |           | git push                                        |
      |           | git checkout prototype                          |
      | prototype | git rebase main --no-update-refs                |
      |           | git push --force-with-lease --force-if-includes |
      |           | git rebase origin/prototype --no-update-refs    |
      |           | git rebase main --no-update-refs                |
      |           | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE       |
      | main      | local, origin | main commit   |
      | prototype | local, origin | origin commit |
      |           |               | local commit  |
    And all branches are now synchronized
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                                     |
      | prototype | git reset --hard {{ sha 'local commit' }}                                                   |
      |           | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin commit' }}:prototype |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE       | FILE NAME   |
      | main      | local, origin | main commit   | main_file   |
      | prototype | local         | local commit  | local_file  |
      |           | origin        | origin commit | origin_file |
    And the initial branches and lineage exist now
