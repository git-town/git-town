@skipWindows
Feature: Azure DevOps support

  Background:
    Given tool "open" is installed

  Scenario Outline: normal origin
    Given the current branch is a feature branch "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://dev.azure.com/organization/repository/compare/feature?expand=1
      """

    Examples:
      | ORIGIN                                        |
      | git@dev.azure.com/organization/repository     |
      | https://dev.azure.com/organization/repository |

  Scenario Outline: origin contains path that looks like a URL
    Given the current branch is a feature branch "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """

    Examples:
      | ORIGIN                                                              |
      | http://dev.azure.com/organization/repository.hosting.azure.com.git  |
      | http://dev.azure.com/organization/repository.hosting.azure.com      |
      | https://dev.azure.com/organization/repository.hosting.azure.com.git |
      | https://dev.azure.com/organization/repository.hosting.azure.com     |
      | git@dev.azure.com:organization/repository.hosting.azure.com.git     |
      | git@dev.azure.com:organization/repository.hosting.azure.com         |

  Scenario Outline: proper URL encoding
    Given the current branch is a feature branch "<BRANCH_NAME>"
    And the origin is "https://dev.azure.com/organization/repository"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      <URL>
      """

    Examples:
      | BRANCH_NAME    | URL                                                                           |
      | feature-branch | https://dev.azure.com/organization/repository/compare/feature-branch?expand=1 |
      | feature_branch | https://dev.azure.com/organization/repository/compare/feature_branch?expand=1 |
      | fix-#2         | https://dev.azure.com/organization/repository/compare/fix-%232?expand=1       |
      | test/feature   | https://dev.azure.com/organization/repository/compare/test%2Ffeature?expand=1 |

  Scenario: nested feature branch with known parent
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the origin is "git@github.com:git-town/git-town.git"
    And the current branch is "child"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent...child?expand=1
      """
