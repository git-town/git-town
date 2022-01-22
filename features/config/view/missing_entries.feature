Feature: "git town config" with unconfigured repo

  To know whether I need to configure a repo
  I want to see missing configuration entries.

  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        [none]

      Perennial branches:
        [none]
      """

  Scenario: the main branch is configured but the perennial branches are not
    Given the main branch is configured as "main"
    And the perennial branches are not configured
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        [none]
      """

  Scenario: the perennial branches are configured but the main branch is not
    Given the main branch name is not configured
    And the perennial branches are configured as "qa" and "staging"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        [none]

      Perennial branches:
        qa
        staging
      """
