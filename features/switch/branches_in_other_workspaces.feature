Feature: switch branches

  Scenario: switching to another branch
    Given a Git repo clone
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
      | gamma | feature | main   | local, origin |
    Given the current branch is "alpha"
    And branch "beta" is active in another worktree
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it runs the commands
      | BRANCH | COMMAND            |
      | alpha  | git checkout gamma |
    And the current branch is now "gamma"
