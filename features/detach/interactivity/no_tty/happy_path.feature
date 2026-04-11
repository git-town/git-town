@skipWindows
Feature: no TTY

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
    And the current branch is "branch-1"
    When I run "git-town detach" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                |
      | branch-1 | git fetch --prune --tags                               |
      |          | git -c rebase.updateRefs=false rebase --onto main main |
