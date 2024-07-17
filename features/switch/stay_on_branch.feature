Feature: stay on the same branch

  Scenario: switching to another branch
    Given a Git repo clone
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
    Given the current branch is "alpha"
    When I run "git-town switch" and enter into the dialogs:
      | KEYS  |
      | enter |
    Then it runs the commands
      | BRANCH | COMMAND |
    And the current branch is still "alpha"
