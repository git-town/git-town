Feature: show help screen even outside of a git repository

  (see ./help_configured.feature)


  Scenario Outline: Running outside of a git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>`
    Then I see the "git-town" man page
    And I don't see "fatal: Not a git repository"

    Examples:
      | COMMAND       |
      | git town      |
      | git town help |
