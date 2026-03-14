Feature: TERM=dumb, missing main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
    And Git Town is not configured
    And the current branch is "feature"
    When I run "git-town compress" with these environment variables
      | TERM | dumb |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      no main branch configured and only a dumb terminal available.

      To configure, run "git config git-town.main-branch <branch>".
      To set up interactively, run "git town init" in a shell with TTY.
      """
