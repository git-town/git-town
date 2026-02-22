@skipWindows
Feature: no TTY

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_1    | content 2    |
    And the current branch is "feature"
    When I run "git-town compress" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git reset --soft main --                        |
      |         | git commit -m "commit 1"                        |
      |         | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 2    |
