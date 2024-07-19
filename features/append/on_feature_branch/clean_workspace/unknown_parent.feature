@messyoutput
Feature: ask for missing parent branch information

  Scenario:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town append new" and enter into the dialog:
      | DIALOG                             | KEYS  |
      | select parent branch for "feature" | enter |
    Then this lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
