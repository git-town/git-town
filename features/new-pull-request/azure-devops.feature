@skipWindows
Feature: Azure DevOps support

  Background:
    Given tool "open" is installed

  @this
  Scenario Outline: normal origin
    Given the current branch is a feature branch "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://dev.azure.com/organization/project/compare/feature?expand=1
      """

    Examples:
      | ORIGIN                                                     |
      | git@ssh.dev.azure.com:v3/organization/project/repository   |
      | https://dev.azure.com/organization/project/_git/repository |

  Scenario Outline: origin contains path that looks like a URL
    Given the current branch is a feature branch "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """

    Examples:
      | ORIGIN                                              |
      | http://github.com/git-town/git-town.github.com.git  |
      | http://github.com/git-town/git-town.github.com      |
      | https://github.com/git-town/git-town.github.com.git |
      | https://github.com/git-town/git-town.github.com     |
      | git@github.com:git-town/git-town.github.com.git     |
      | git@github.com:git-town/git-town.github.com         |

  Scenario Outline: proper URL encoding
    Given the current branch is a feature branch "<BRANCH_NAME>"
    And the origin is "https://github.com/git-town/git-town"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      <URL>
      """

    Examples:
      | BRANCH_NAME    | URL                                                                  |
      | feature-branch | https://github.com/git-town/git-town/compare/feature-branch?expand=1 |
      | feature_branch | https://github.com/git-town/git-town/compare/feature_branch?expand=1 |
      | fix-#2         | https://github.com/git-town/git-town/compare/fix-%232?expand=1       |
      | test/feature   | https://github.com/git-town/git-town/compare/test%2Ffeature?expand=1 |

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
