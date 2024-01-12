@skipWindows
Feature: ask for missing parent branch information

  @debug @this
  Scenario:
    Given the current branch is "feature"
    And inspect the repo
    When I run "git-town append new" and enter into the dialog:
      | KEY   | DESCRIPTION        |
      | enter | select this option |
    Then this branch lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
