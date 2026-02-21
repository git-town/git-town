@skipWindows
Feature: non-TTY usage

  Scenario Outline:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | branch-1  | (none)  |           | local, origin |
      | feature-1 | feature | main      | local, origin |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
    And Git Town is not configured
    And the current branch is "<BRANCH>"
    When I run "git-town <COMMAND>" in a non-TTY shell
    Then Git Town prints the error:
      """
      no main branch configured and no interactive terminal available.
      To configure, run "git config git-town.main-branch <branch>".
      To set up interactively, run "git town init" in a shell with TTY.
      """

    Examples:
      | BRANCH   | COMMAND                        |
      | branch-1 | append new                     |
      | branch-1 | commit --down --message commit |
      | branch-1 | compress                       |
      | branch-1 | delete                         |
      | branch-1 | detach                         |
      | branch-1 | diff-parent                    |
      | branch-1 | hack new                       |
      | branch-1 | merge                          |
      | branch-1 | prepend new                    |
      | branch-1 | propose                        |
      | branch-1 | rename new                     |
      | branch-1 | set-parent new                 |
      | branch-1 | ship                           |
      | branch-1 | skip                           |
      | branch-1 | swap                           |
      | branch-1 | sync                           |
      | branch-1 | undo                           |
      | branch-1 | walk --all                     |
