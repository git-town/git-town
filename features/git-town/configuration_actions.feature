Feature: Configuration actions

  Scenario: Reading configuration
    Given I have set "main-branch-name" to "main"
    And I have set "non-feature-branch-names" to "dev, qa, staging"
    When I run `git town --config` while allowing errors
    Then the output should contain "Main branch: 'main'"
    And the output should contain "Non-feature branches: 'dev, qa, staging'"
