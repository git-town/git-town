@skipWindows
Feature: GitLab support

  Background:
    Given tool "open" is installed

  Scenario Outline: creating pull-requests
    Given the current branch is a feature branch "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitlab.com/kadu/kadu/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main
      """

    Examples:
      | ORIGIN                           |
      | https://gitlab.com/kadu/kadu.git |
      | git@gitlab.com:kadu/kadu.git     |

  Scenario: nested feature branch with known parent
    Given a feature branch "parent-feature"
    And a feature branch "child-feature" as a child of "parent-feature"
    And the origin is "git@gitlab.com:kadu/kadu.git"
    And the current branch is "child-feature"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitlab.com/kadu/kadu/merge_requests/new?merge_request%5Bsource_branch%5D=child-feature&merge_request%5Btarget_branch%5D=parent-feature
      """
