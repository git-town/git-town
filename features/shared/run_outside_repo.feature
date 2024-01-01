Feature: require a Git repository

  Scenario Outline:
    Given I am outside a Git repo
    When I run "<COMMAND>"
    Then it prints the error:
      """
      this is not a Git repository
      """

    Examples:
      | COMMAND                            |
      | git-town config                    |
      | git-town config main-branch        |
      | git-town config perennial-branches |
      | git-town hack feature              |
      | git-town kill                      |
      | git-town propose                   |
      | git-town repo                      |
      | git-town ship                      |
      | git-town sync                      |
