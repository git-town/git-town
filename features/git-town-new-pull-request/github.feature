Feature: git-new-pull-request when origin is on GitHub

  As a developer having finished a feature in a repository hosted on GitHub
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Background:
    Given I have "open" installed


  Scenario Outline: creating pull-requests
    Given I have a feature branch named "feature"
    And my remote origin is <ORIGIN>
    And I am on the "feature" branch
    When I run `gt new-pull-request`
    Then I see a new GitHub pull request for the "feature" branch in the "<REPOSITORY>" repo in my browser

    Examples:
      | ORIGIN                                                | REPOSITORY                     |
      | http://github.com/Originate/git-town.git              | Originate/git-town             |
      | http://github.com/Originate/git-town                  | Originate/git-town             |
      | https://github.com/Originate/git-town.git             | Originate/git-town             |
      | https://github.com/Originate/git-town                 | Originate/git-town             |
      | git@github.com:Originate/git-town.git                 | Originate/git-town             |
      | git@github.com:Originate/git-town                     | Originate/git-town             |
      | git@github-as-account1:Originate/git-town.git         | Originate/git-town             |
      | http://github.com/Originate/originate.github.com.git  | Originate/originate.github.com |
      | http://github.com/Originate/originate.github.com      | Originate/originate.github.com |
      | https://github.com/Originate/originate.github.com.git | Originate/originate.github.com |
      | https://github.com/Originate/originate.github.com     | Originate/originate.github.com |
      | git@github.com:Originate/originate.github.com.git     | Originate/originate.github.com |
      | git@github.com:Originate/originate.github.com         | Originate/originate.github.com |
      | ssh://git@github.com/Originate/git-town.git           | Originate/git-town             |
      | ssh://git@github.com/Originate/git-town               | Originate/git-town             |


  Scenario: nested feature branch with known parent
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"
    And my remote origin is git@github.com:Originate/git-town.git
    And I am on the "child-feature" branch
    When I run `gt new-pull-request`
    Then I see a new GitHub pull request for the "child-feature" branch against the "parent-feature" branch in the "Originate/git-town" repo in my browser


  Scenario: nested feature branch with unknown parent (entering the parent name)
    Given I have a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And my remote origin is git@github.com:Originate/git-town.git
    And I am on the "feature" branch
    When I run `gt new-pull-request` and enter "main"
    Then I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser


  Scenario: nested feature branch with unknown parent (accepting default choice)
    Given I have a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And my remote origin is git@github.com:Originate/git-town.git
    And I am on the "feature" branch
    When I run `gt new-pull-request` and press ENTER
    Then I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser
