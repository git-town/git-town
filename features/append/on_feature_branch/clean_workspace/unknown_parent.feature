@messyoutput
Feature: ask for missing parent branch information

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT |
      | feature | feature | main   |
    And the current branch is "feature"
    When I run "git-town append new" and enter into the dialog:
      | DIALOG                             | KEYS  |
      | select parent branch for "feature" | enter |
    Then this lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
