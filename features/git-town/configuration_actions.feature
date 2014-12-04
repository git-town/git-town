Feature: listing the configuration

  Scenario: Everything is configured
    Given I have configured the main branch name as "main"
    And my non-feature branches are configured as "dev, qa, staging"
    When I run `git town --config`
    Then I see "Main branch: 'main'"
    And I see "Non-feature branches: 'dev, qa, staging'"
