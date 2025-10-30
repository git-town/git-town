Feature: compress the branch that has the same name as a folder

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE  | FILE NAME |
      | main   | local    | commit 1 | main      |
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local    | commit 1 | file_1    | content 1    |
      |         |          | commit 2 | file_2    | content 2    |
    And the current branch is "feature"
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git reset --soft main -- |
      |         | git commit -m "commit 1" |
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE  |
      | main    | local    | commit 1 |
      | feature | local    | commit 1 |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                               |
      | feature | git reset --hard {{ sha 'commit 2' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
