Feature: show the configuration

  Scenario: all configured, no nested branches
    Given the main branch is "main"
    And the perennial branches are "qa" and "staging"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        qa
        staging
      """

  Scenario: all configured, with nested branches
    Given the main branch is "main"
    And my repo has the perennial branches "qa" and "staging"
    And my repo has the feature branches "alpha" and "beta"
    And my repo has a feature branch "child" as a child of "alpha"
    And my repo has a feature branch "hotfix" as a child of "qa"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        qa
        staging

      Branch Ancestry:
        main
          alpha
            child
          beta

        qa
          hotfix
      """

  Scenario: no configuration data
    Given I haven't configured Git Town yet
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        [none]

      Perennial branches:
        [none]
      """
