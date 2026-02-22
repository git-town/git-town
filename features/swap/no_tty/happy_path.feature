@skipWindows
Feature: no TTY

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE  |
      | branch-1 | local    | commit 1 |
      | branch-2 | local    | commit 2 |
    And the current branch is "branch-2"
    When I run "git-town swap" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git checkout branch-1                                      |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto branch-2 main |
      |          | git checkout branch-2                                      |
    And this lineage exists now
      """
      main
        branch-2
          branch-1
      """

  Scenario: undo
    When I run "git-town undo" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH   | COMMAND                               |
      | branch-2 | git checkout branch-1                 |
      | branch-1 | git reset --hard {{ sha 'commit 1' }} |
      |          | git checkout branch-2                 |
    And this lineage exists now
      """
      main
        branch-1
          branch-2
      """
