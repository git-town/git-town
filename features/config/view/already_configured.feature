Feature: listing the configuration

  To manage Git Town's configuration efficiently
  I want to see the complete Git Town configuration with one command.

  Scenario: without nested branches
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

  Scenario: with nested branches
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

      Pull branch strategy:
        rebase

      New Branch Push Flag:
        false
      """
