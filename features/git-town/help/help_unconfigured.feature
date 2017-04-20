Feature: show help screen when Git Town is not configured

  (see ./help_configured.feature)


  Background:
    Given I haven't configured Git Town yet


  Scenario Outline:
    When I run `<COMMAND>`
    Then I see
      """
      Usage:
        gt [command]
      """

    Examples:
      | COMMAND |
      | gt      |
      | gt help |
