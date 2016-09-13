Feature: show help screen even outside of a git town-repository

  (see ./help_configured.feature)


  Scenario Outline: Running outside of a git town-repository
    Given I'm currently not in a git town-repository
    When I run `<COMMAND>`
    Then I see the "git-town" man page
    And I don't see "fatal: Not a git town-repository"

    Examples:
      | COMMAND       |
      | git town      |
      | git town help |
