Feature: git-new-pull-request when origin is on GitLab

  As a developer having finished a feature in a repository hosted on GitLab
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Background:
    Given I have "open" installed


  Scenario Outline: creating pull-requests
    Given I have a feature branch named "feature"
    And my remote origin is <ORIGIN>
    And I am on the "feature" branch
    When I run `git town-new-pull-request`
    Then I see a new GitLab pull request for the "feature" branch in the "<REPOSITORY>" repo in my browser

    Examples:
      | ORIGIN                           | REPOSITORY |
      | https://gitlab.com/kadu/kadu.git | kadu/kadu  |
      | git@gitlab.com:kadu/kadu.git     | kadu/kadu  |


  Scenario: nested feature branch with known parent
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"
    And my remote origin is git@gitlab.com:kadu/kadu.git
    And I am on the "child-feature" branch
    When I run `git town-new-pull-request`
    Then I see a new GitLab pull request for the "child-feature" branch against the "parent-feature" branch in the "kadu/kadu" repo in my browser


  Scenario: nested feature branch with unknown parent (entering the parent name)
    Given I have a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And my remote origin is git@gitlab.com:kadu/kadu.git
    And I am on the "feature" branch
    When I run `git town-new-pull-request` and enter "main"
    Then I see a new GitLab pull request for the "feature" branch in the "kadu/kadu" repo in my browser


  Scenario: nested feature branch with unknown parent (accepting default choice)
    Given I have a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And my remote origin is git@gitlab.com:kadu/kadu.git
    And I am on the "feature" branch
    When I run `git town-new-pull-request` and press ENTER
    Then I see a new GitLab pull request for the "feature" branch in the "kadu/kadu" repo in my browser
