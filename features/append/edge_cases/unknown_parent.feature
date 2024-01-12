Feature: ask for missing parent branch information

  @debug
  Scenario:
    Given the current branch is "feature"
    When I run "git-town append new" and enter into the dialog:
      | DIALOG             | KEYS  |
      | select main branch | enter |
    Then this branch lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
