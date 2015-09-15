Feature: Initial configuration

  As a user who hasn't configured Git Town yet
  I want to have a simple, dedicated setup command
  So that I can configure it safely before using any Git Town command


  Background:
    Given I haven't configured Git Town yet


  Scenario: user enters valid main branch and perennial branches
    Given I have branches named "qa" and "production"
    When I run `git hack feature` and enter "main", "qa", "production" and ""
    Then the main branch name is now configured as "main"
    And my perennial branches are now configured as "qa" and "production"


  Scenario: user enters empty main branch
    When I run `git hack feature` and enter "", "main" and ""
    Then I see "no input received"
    And the main branch name is now configured as "main"
    And my perennial branches are configured as none


  Scenario: user enters invalid main branch number
    When I run `git hack feature` and enter "2", "main" and ""
    Then I see "Invalid branch number"
    And the main branch name is now configured as "main"
    And my perennial branches are configured as none


  Scenario: user enters non-existent main branch
    When I run `git hack feature` and enter "non-existent", "main" and ""
    Then I see "branch 'non-existent' doesn't exist"
    And the main branch name is now configured as "main"
    And my perennial branches are configured as none


  Scenario: user enters main branch as perennial branch
    When I run `git hack feature` and enter "main", "main" and ""
    Then I see "'main' is already set as the main branch"
    And the main branch name is now configured as "main"
    And my perennial branches are configured as none


  Scenario: user enters invalid perennial branch number
    When I run `git hack feature` and enter "main", "2" and ""
    Then I see "Invalid branch number"
    And the main branch name is now configured as "main"
    And my perennial branches are configured as none


  Scenario: user enters non-existent perennial branch
    When I run `git hack feature` and enter "main", "non-existent" and ""
    Then I see "branch 'non-existent' doesn't exist"
    And the main branch name is now configured as "main"
    And my perennial branches are configured as none
