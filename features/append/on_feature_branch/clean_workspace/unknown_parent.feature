Feature: ask for missing parent branch information

  Scenario:
    Given the current branch is "feature"
    When I run "git-town append new" and enter into the dialog:
      | DIALOG             | KEYS  |
      | select main branch | enter |
    Then this lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
