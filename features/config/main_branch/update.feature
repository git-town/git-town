Feature: configure the main branch

  Scenario: not configured, no config file
    When I run "git-town config main-branch main"
    Then it prints no output
    And local Git Town setting "main-branch" is now "main"
    And still no configuration file exists

  Scenario: empty local Git setting
    Given local Git Town setting "main-branch" is ""
    When I run "git-town config main-branch main"
    Then it prints:
      """
      NOTICE: cleaned up empty configuration entry "git-town.main-branch"
      """
    And local Git Town setting "main-branch" is now "main"
    And still no configuration file exists

  Scenario: update existing local Git setting
    Given the branches "old" and "new"
    And local Git Town setting "main-branch" is "old"
    When I run "git-town config main-branch new"
    Then it prints no output
    And local Git Town setting "main-branch" is now "new"
    And still no configuration file exists

  Scenario: update to non-existing branch
    When I run "git-town config main-branch non-existing"
    Then it prints the error:
      """
      there is no branch "non-existing"
      """

  Scenario: not configured, config file exists
    Given local Git Town setting "main-branch" doesn't exist
    And a branch "new"
    And the configuration file:
      """
      """
    When I run "git-town config main-branch new"
    Then it prints no output
    And the configuration file is now:
      """
      [branches]
        main = "new"
      """
    And local Git Town setting "main-branch" still doesn't exist

  Scenario: existing entry in config file
    Given local Git Town setting "main-branch" doesn't exist
    And a branch "new"
    And the configuration file:
      """
      [branches]
        main = "old"
      """
    When I run "git-town config main-branch new"
    Then it prints no output
    And the configuration file is now:
      """
      [branches]
        main = "new"
      """
    And local Git Town setting "main-branch" still doesn't exist
