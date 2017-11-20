Feature: git-new-pull-request when origin is on GitLab

  As a developer having finished a feature in a repository hosted on GitLab
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Background:
    Given I have "open" installed


  Scenario Outline: creating pull-requests
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is <ORIGIN>
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      https://gitlab.com/kadu/kadu/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main
      """

    Examples:
      | ORIGIN                           |
      | https://gitlab.com/kadu/kadu.git |
      | git@gitlab.com:kadu/kadu.git     |


  Scenario: nested feature branch with known parent
    Given my repository has a feature branch named "parent-feature"
    And my repository has a feature branch named "child-feature" as a child of "parent-feature"
    And my repo's remote origin is git@gitlab.com:kadu/kadu.git
    And I am on the "child-feature" branch
    When I run `git-town new-pull-request`
    Then I see a new GitLab pull request for the "child-feature" branch against the "parent-feature" branch in the "kadu/kadu" repo in my browser


  Scenario: nested feature branch with unknown parent (entering the parent name)
    Given my repository has a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And my repo's remote origin is git@gitlab.com:kadu/kadu.git
    And I am on the "feature" branch
    When I run `git-town new-pull-request` and enter "main"
    Then I see a new GitLab pull request for the "feature" branch in the "kadu/kadu" repo in my browser


  Scenario: nested feature branch with unknown parent (accepting default choice)
    Given my repository has a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And my repo's remote origin is git@gitlab.com:kadu/kadu.git
    And I am on the "feature" branch
    When I run `git-town new-pull-request` and press ENTER
    Then I see a new GitLab pull request for the "feature" branch in the "kadu/kadu" repo in my browser
