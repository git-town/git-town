Feature: show help screen even outside of a Git repository

  (see ./help_configured.feature)


  Scenario Outline: Running outside of a Git repository
    Given I'm currently not in a Git repository
    When I run `<COMMAND>`
    Then I see the "git-town" man page
    And I don't see "fatal: Not a Git repository"

    Examples:
      | COMMAND       |
      | git town      |
      | git town help |
