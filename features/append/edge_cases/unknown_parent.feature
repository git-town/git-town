@skipWindows
Feature: ask for missing parent branch information

  @this
  Scenario:
    Given the current branch is "feature"
    When I run "git-town append new" and enter into the dialog:
      | KEY   | DESCRIPTION        |
      | enter | select main branch |
    Then this branch lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
