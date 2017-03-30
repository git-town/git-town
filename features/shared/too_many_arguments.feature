Feature: too many arguments

  As a developer providing too many arguments
  I should be reminded that the command expects no arguments
  So that I can use it correctly without having to look that fact up in the readme.


  Scenario: hack
    When I run `gt hack feature1`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt hack <branch> [flags]
      """


  Scenario: kill
    When I run `gt kill feature1 feature2`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt kill [<branch>] [flags]
      """


  Scenario: prune-branches
    When I run `gt prune-branches feature1`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt prune-branches [flags]
      """


  Scenario: sync
    When I run `gt sync feature1`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt sync [flags]
      """
