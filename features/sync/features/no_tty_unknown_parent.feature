Feature: descriptive error when parent is unknown and no TTY is available

  Scenario: sync a branch with unknown parent and no TTY
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE   | LOCATIONS     |
      | alpha | (none) | local, origin |
    And the current branch is "alpha"
    When I run "git-town sync"
    Then Git Town prints the error:
      """
      no parent configured for branch "alpha" and no interactive terminal available.
      To set the parent, run: git town set-parent <parent>
      To configure manually, run: git config git-town-branch.alpha.parent <parent>
      """
