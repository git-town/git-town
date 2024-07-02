@messyoutput
Feature: ask for missing parent branch information

  Scenario:
    Given the branches
      | NAME    | TYPE    |
      | main    | main    |
      | feature | feature |
    And the current branch is "feature"
    When I run "git-town append new" and enter into the dialog:
      | DIALOG                             | KEYS  |
      | select parent branch for "feature" | enter |
    Then this lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
