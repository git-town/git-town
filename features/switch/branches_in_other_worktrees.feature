Feature: switch while some branches are checked out in another worktree

  @this
  Scenario: some branches are checked out in another worktree
    Given the current branch is a local feature branch "alpha"
    And a local feature branch "beta"
    And a local feature branch "gamma"
    And branch "beta" is active in another worktree
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it runs the commands
      | BRANCH | COMMAND            |
      | alpha  | git checkout gamma |
    And the current branch is now "gamma"
