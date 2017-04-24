Feature: listing the configuration

  As a user unsure about how Git Town is currently configured
  I want to be able to see the complete Git Town configuration with one command
  So that I can configure Git Town efficiently and have more time for actual work.


  Scenario: everything is configured
    Given I have configured the main branch name as "main"
    And my perennial branches are configured as "qa" and "staging"
    When I run `gt config`
    Then I see
      """
      Main branch:
        main

      Perennial branches:
        qa
        staging
      """


  Scenario: everything is configured and there are nested branches
    Given I have configured the main branch name as "main"
    And my perennial branches are configured as "qa" and "staging"
    And I have feature branches named "parent-feature" and "stand-alone-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"
    When I run `gt config`
    Then I see
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
      """


  Scenario: the main branch is configured but the perennial branches are not
    Given I have configured the main branch name as "main"
    And my perennial branches are not configured
    When I run `gt config`
    Then I see
      """
      Main branch:
        main

      Perennial branches:
        [none]
      """


  Scenario: the main branch is not configured but the perennial branches are
    Given I don't have a main branch name configured
    And my perennial branches are configured as "qa" and "staging"
    When I run `gt config`
    Then I see
      """
      Main branch:
        [none]

      Perennial branches:
        qa
        staging
      """


  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run `gt config`
    Then I see
      """
      Main branch:
        [none]

      Perennial branches:
        [none]
      """
