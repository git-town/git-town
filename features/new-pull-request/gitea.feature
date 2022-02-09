@skipWindows
Feature: Gitea support

  Background:
    Given the "open" tool is installed

  Scenario Outline: normal origin
    Given the current branch is a feature branch "feature"
    And my repo's origin is "<ORIGIN>"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitea.com/git-town/git-town/compare/main...feature
      """

    Examples:
      | ORIGIN                                    |
      | http://gitea.com/git-town/git-town.git    |
      | http://gitea.com/git-town/git-town        |
      | https://gitea.com/git-town/git-town.git   |
      | https://gitea.com/git-town/git-town       |
      | git@gitea.com:git-town/git-town.git       |
      | git@gitea.com:git-town/git-town           |
      | ssh://git@gitea.com/git-town/git-town.git |
      | ssh://git@gitea.com/git-town/git-town     |

  Scenario Outline: origin contains path that looks like a URL
    Given the current branch is a feature branch "feature"
    And my repo's origin is "<ORIGIN>"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitea.com/git-town/git-town.gitea.com/compare/main...feature
      """

    Examples:
      | ORIGIN                                            |
      | http://gitea.com/git-town/git-town.gitea.com.git  |
      | http://gitea.com/git-town/git-town.gitea.com      |
      | https://gitea.com/git-town/git-town.gitea.com.git |
      | https://gitea.com/git-town/git-town.gitea.com     |
      | git@gitea.com:git-town/git-town.gitea.com.git     |
      | git@gitea.com:git-town/git-town.gitea.com         |

  Scenario Outline: proper URL encoding
    Given the current branch is a feature branch "<BRANCH_NAME>"
    And my repo's origin is "https://gitea.com/git-town/git-town"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      <URL>
      """

    Examples:
      | BRANCH_NAME    | URL                                                               |
      | feature-branch | https://gitea.com/git-town/git-town/compare/main...feature-branch |
      | feature_branch | https://gitea.com/git-town/git-town/compare/main...feature_branch |
      | fix-#2         | https://gitea.com/git-town/git-town/compare/main...fix-%232       |
      | test/feature   | https://gitea.com/git-town/git-town/compare/main...test%2Ffeature |

  Scenario: nested feature branch with known parent
    Given a feature branch "parent-feature"
    And a feature branch "child-feature" as a child of "parent-feature"
    And my repo's origin is "git@gitea.com:git-town/git-town.git"
    And the current branch is "child-feature"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitea.com/git-town/git-town/compare/parent-feature...child-feature
      """
