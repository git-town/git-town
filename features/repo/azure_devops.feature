@skipWindows
Feature: Azure DevOps

  Scenario Outline:
    Given the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://dev.azure.com/organization/repository
      """

    Examples:
      | ORIGIN                                                     |
      | http://username@dev.azure.com/organization/repository.git  |
      | http://username@dev.azure.com/organization/repository      |
      | https://username@dev.azure.com/organization/repository.git |
      | https://username@dev.azure.com/organization/repository     |
      | git@dev.azure.com/organization/repository.git              |
      | git@dev.azure.com/organization/repository                  |
      | ssh://git@dev.azure.com/organization/repository.git        |
      | ssh://git@dev.azure.com/organization/repository            |
      | username@dev.azure.com/organization/repository.git         |
      | username@dev.azure.com/organization/repository             |
      | ssh://username@dev.azure.com/organization/repository.git   |
      | ssh://username@dev.azure.com/organization/repository       |
