Feature: switch branches

  Scenario: switching to another branch while merging open changes
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the current branch is "current"
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And an uncommitted file
    When I run "git-town switch -m" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | current | git checkout other -m |
    And the current branch is now "other"
