Feature: git-new-pull-request when origin is on GitHub

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
      https://github.com/git-town/git-town/compare/feature?expand=1
      """

    Examples:
      | ORIGIN                                   |
      | http://github.com/git-town/git-town.git  |
      | http://github.com/git-town/git-town      |
      | https://github.com/git-town/git-town.git |
      | https://github.com/git-town/git-town     |
      | git@github.com:git-town/git-town.git     |
      | git@github.com:git-town/git-town         |

  @skipWindows
  Scenario Outline: origin contains path that looks like a URL
    Given my repo has a feature branch named "feature"
    And my repo's origin is "<ORIGIN>"
    And I am on the "feature" branch
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

  @skipWindows
  Scenario Outline: proper URL encoding
    Given my repo has a feature branch named "<BRANCH_NAME>"
    And my repo's origin is "https://github.com/git-town/git-town"
    And I am on the "<BRANCH_NAME>" branch
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

  @skipWindows
  Scenario Outline: SSH style origin
    Given my repo has a feature branch named "feature"
    And my repo's origin is "<ORIGIN>"
    And I am on the "feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """

    Examples:
      | ORIGIN                                     |
      | ssh://git@github.com/git-town/git-town.git |
      | ssh://git@github.com/git-town/git-town     |

  @skipWindows
  Scenario: nested feature branch with known parent
    Given my repo has a feature branch named "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And my repo's origin is "git@github.com:git-town/git-town.git"
    And I am on the "child-feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent-feature...child-feature?expand=1
      """
