Feature: listing the configuration

  Scenario: all configured, no nested branches
    Given the main branch is configured as "main"
    And the perennial branches are configured as "qa" and "staging"
    When I run "git-town config"
    Then it prints:
      """
      Main branch:
        main

      Perennial branches:
        qa
        staging
      """

  Scenario: all configured, has nested branches
    Given the main branch is configured as "main"
    And my repo has the perennial branches "qa" and "staging"
    And my repo has the feature branches "parent-feature" and "stand-alone-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And my repo has a feature branch named "qa-hotfix" as a child of "qa"
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
          parent-feature
            child-feature
          stand-alone-feature

        qa
          qa-hotfix
      """

  Scenario: main branch is configured, perennial branches are not
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

  Scenario: perennial branches are configured, main branch is not
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
