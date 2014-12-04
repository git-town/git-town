Feature: Configuration

  Scenario: Without a configured main branch name
    Given I don't have a main branch name configured
    When I run `git ship` and enter "user_main_branch"
    Then the main branch name is now configured as "user_main_branch"


  Scenario: Automatic update of old configuration files
    Given I have an old configuration file with main branch: "main"
    When I run `git ship` while allowing errors
    Then the main branch name is now configured as "main"
    And I don't have an old configuration file anymore


  Scenario: Reading configuration
    Given I have set "main-branch-name" to "main"
    And I have set "non-feature-branch-names" to "dev, qa, staging"
    When I run `git town --config` while allowing errors
    Then the output should contain 'Main branch: "main"'
    And the output should contain 'Non-feature branches: "dev, qa, staging"'
