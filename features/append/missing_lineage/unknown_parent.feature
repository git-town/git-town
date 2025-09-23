@messyoutput
Feature: ask for missing parent branch information

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | branch | (none) |        | local, origin |
    And the current branch is "branch"
    When I run "git-town append new" and enter into the dialog:
      | DIALOG                     | KEYS  |
      | parent branch for "branch" | enter |
    Then this lineage exists now
      """
      main
        branch
          new
      """
