Feature: switch branches

  Scenario: switching to another branch with conflicting open changes
    Given the current branch is a feature branch "current"
    And a feature branch "other"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | other  | local, origin | other commit |
    And an uncommitted file
    When I run "git-town switch -m" and enter into the dialogs:
      | DIALOG | KEYS       |
      |        | down enter |
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | current | git checkout other -m |
    And the current branch is now "other"
