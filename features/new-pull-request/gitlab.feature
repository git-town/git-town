Feature: git-new-pull-request when origin is on GitLab

  As a developer having finished a feature in a repository hosted on GitLab
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Background:
    Given my computer has the "open" tool installed


  Scenario Outline: creating pull-requests
    Given my repo has a feature branch named "feature"
    And my repo's origin is "<ORIGIN>"
    And I am on the "feature" branch
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
    Given my repo has a feature branch named "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And my repo's origin is "git@gitlab.com:kadu/kadu.git"
    And I am on the "child-feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitlab.com/kadu/kadu/merge_requests/new?merge_request%5Bsource_branch%5D=child-feature&merge_request%5Btarget_branch%5D=parent-feature
      """
