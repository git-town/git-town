Feature: git-new-pull-request when origin is on Gitea

  Background:
    Given my computer has the "open" tool installed

  @skipWindows
  Scenario Outline: normal origin
    Given my repo has a feature branch named "feature"
    And my repo's origin is "<ORIGIN>"
    And I am on the "feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitea.com/git-town/git-town/compare/main...feature
      """

    Examples:
      | ORIGIN                                  |
      | http://gitea.com/git-town/git-town.git  |
      | http://gitea.com/git-town/git-town      |
      | https://gitea.com/git-town/git-town.git |
      | https://gitea.com/git-town/git-town     |
      | git@gitea.com:git-town/git-town.git     |
      | git@gitea.com:git-town/git-town         |

  @skipWindows
  Scenario Outline: origin contains path that looks like a URL
    Given my repo has a feature branch named "feature"
    And my repo's origin is "<ORIGIN>"
    And I am on the "feature" branch
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

  @skipWindows
  Scenario Outline: proper URL encoding
    Given my repo has a feature branch named "<BRANCH_NAME>"
    And my repo's origin is "https://gitea.com/git-town/git-town"
    And I am on the "<BRANCH_NAME>" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      <URL>
      """

    Examples:
      | BRANCH_NAME    | URL                                                        |
      | feature-branch | https://gitea.com/git-town/git-town/compare/main...feature-branch |
      | feature_branch | https://gitea.com/git-town/git-town/compare/main...feature_branch |
      | fix-#2         | https://gitea.com/git-town/git-town/compare/main...fix-%232       |
      | test/feature   | https://gitea.com/git-town/git-town/compare/main...test%2Ffeature |

  @skipWindows
  Scenario Outline: SSH style origin
    Given my repo has a feature branch named "feature"
    And my repo's origin is "<ORIGIN>"
    And I am on the "feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitea.com/git-town/git-town/compare/main...feature
      """

    Examples:
      | ORIGIN                                    |
      | ssh://git@gitea.com/git-town/git-town.git |
      | ssh://git@gitea.com/git-town/git-town     |

  @skipWindows
  Scenario: nested feature branch with known parent
    Given my repo has a feature branch named "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And my repo's origin is "git@gitea.com:git-town/git-town.git"
    And I am on the "child-feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitea.com/git-town/git-town/compare/parent-feature...child-feature
      """
