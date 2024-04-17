Feature: switch branches

  @this
  Scenario: switching to another branch
    Given the current branch is a feature branch "alpha"
    And a feature branch "beta"
    And a feature branch "gamma"
    And branch "beta" is active in another worktree
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it runs the commands
      | BRANCH | COMMAND            |
      | alpha  | git checkout gamma |
    And the current branch is now "gamma"
