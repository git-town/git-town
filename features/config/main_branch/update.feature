@this
Feature: configure the main branch

  Scenario: not configured
    When I run "git-town config main-branch main"
    Then it prints no output
    And local Git Town setting "main-branch" is now "main"

  Scenario: empty local Git setting
    Given local Git Town setting "main-branch" is ""
    When I run "git-town config main-branch main"
    Then it prints:
      """
      NOTICE: cleaned up empty configuration entry "git-town.main-branch"
      """
    And local Git Town setting "main-branch" is now "main"

  Scenario: update existing local Git setting
    Given the branches "old" and "new"
    And local Git Town setting "main-branch" is "old"
    When I run "git-town config main-branch new"
    Then it prints no output
    And local Git Town setting "main-branch" is now "new"

  Scenario: update to non-existing branch
    When I run "git-town config main-branch non-existing"
    Then it prints the error:
      """
      there is no branch "non-existing"
      """
