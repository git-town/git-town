Feature: git-new-pull-request when origin is on GitHub

  As a developer having finished a feature in a repository hosted on GitHub
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Background:
    Given I have "open" installed


  Scenario Outline: normal origin
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is <ORIGIN>
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
		Then I see a new pull request with this url in my browser:
		  """
			https://github.com/Originate/git-town/compare/feature?expand=1
			"""

    Examples:
      | ORIGIN                                    |
      | http://github.com/Originate/git-town.git  |
      | http://github.com/Originate/git-town      |
      | https://github.com/Originate/git-town.git |
      | https://github.com/Originate/git-town     |
      | git@github.com:Originate/git-town.git     |
      | git@github.com:Originate/git-town         |


  Scenario Outline: origin contains path that looks like a URL
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is <ORIGIN>
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
		Then I see a new pull request with this url in my browser:
		  """
			https://github.com/Originate/originate.github.com/compare/feature?expand=1 |
			"""

		Examples:
      | ORIGIN                                                |
      | http://github.com/Originate/originate.github.com.git  |
      | http://github.com/Originate/originate.github.com      |
      | https://github.com/Originate/originate.github.com.git |
      | https://github.com/Originate/originate.github.com     |
      | git@github.com:Originate/originate.github.com.git     |
      | git@github.com:Originate/originate.github.com         |


  Scenario Outline: proper URL encoding
    Given my repository has a feature branch named "<BRANCH_NAME>"
    And my repo's remote origin is "https://github.com/Originate/git-town"
    And I am on the "<BRANCH_NAME>" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      <URL>
      """

    Examples:
      | BRANCH_NAME    | URL                                                                   |
      | feature-branch | https://github.com/Originate/git-town/compare/feature-branch?expand=1 |
      | feature_branch | https://github.com/Originate/git-town/compare/feature_branch?expand=1 |
      | fix-#2         | https://github.com/Originate/git-town/compare/fix-%232?expand=1       |
      | test/feature   | https://github.com/Originate/git-town/compare/test%2Ffeature?expand=1 |


  Scenario Outline: SSH style origin
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is <ORIGIN>
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
		Then I see a new pull request with this url in my browser:
		  """
			https://github.com/Originate/git-town/compare/feature?expand=1
			"""

		Examples:
      | ORIGIN                                      |
      | ssh://git@github.com/Originate/git-town.git |
      | ssh://git@github.com/Originate/git-town     |


  Scenario: nested feature branch with known parent
    Given my repository has a feature branch named "parent-feature"
    And my repository has a feature branch named "child-feature" as a child of "parent-feature"
    And my repo's remote origin is git@github.com:Originate/git-town.git
    And I am on the "child-feature" branch
    When I run `git-town new-pull-request`
    Then I see a new GitHub pull request for the "child-feature" branch against the "parent-feature" branch in the "Originate/git-town" repo in my browser
