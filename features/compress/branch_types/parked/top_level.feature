Feature: compress the commits on a parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | parked | local, origin | commit 1 | file_1    | content 1    |
      |        |               | commit 2 | file_2    | content 2    |
      |        |               | commit 3 | file_3    | content 3    |
    And the current branch is "parked"
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | parked | git fetch --prune --tags                        |
      |        | git reset --soft main --                        |
      |        | git commit -m "commit 1"                        |
      |        | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | parked | local, origin | commit 1 |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | parked | git reset --hard {{ sha 'commit 3' }}           |
      |        | git push --force-with-lease --force-if-includes |
    And the initial branches and lineage exist now
    And the initial commits exist now
