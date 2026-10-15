@messyoutput
Feature: select a branch that is active in another worktree as the missing parent

  Scenario: feature branch without parent
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
      | branch-2 | (none)  |        | local     |
    And the current branch is "branch-2"
    And branch "branch-1" is active in another worktree
    When I run "git-town diff-parent" and enter into the dialog:
      | DIALOG                       | KEYS       |
      | parent branch for "branch-2" | down enter |
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                 |
      | branch-2 | git diff --merge-base branch-1 branch-2 |
    And this lineage exists now
      """
      main
        branch-1
          branch-2
      """
