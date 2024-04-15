Feature: switch branches

  # @debug @this
  Scenario: switching to another branch
    Given the current branch is a feature branch "alpha"
    And a feature branch "beta"
    When I run "git-town switch" and enter into the dialogs:
      | DIALOG  | KEYS       |
      | welcome | down enter |
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
      |         | git push                 |
      |         | git checkout feature     |
    And I am now on the "beta" branch
