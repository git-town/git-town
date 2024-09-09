Feature: switch to a new remote branch

  @this
  Scenario: switching to another branch while merging open changes
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | local-1  | feature | main   | local, origin |
      | local-2  | feature | main   | local, origin |
      | remote-1 | feature | main   | origin        |
    And the current branch is "local-2"
    And an uncommitted file
    When I run "git-town switch --all" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | current | git checkout remote-1 |
    And the current branch is now "remote-1"
