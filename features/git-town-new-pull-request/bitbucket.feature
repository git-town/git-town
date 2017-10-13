Feature: git-new-pull-request when origin is on Bitbucket

  As a developer having finished a feature in a repository hosted on Bitbucket
  I want to be able to quickly create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Scenario Outline: normal origin
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is <ORIGIN>
    And I have "open" installed
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
		Then I see a new pull request with this url in my browser:
		  """
			https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature
			"""

    Examples:
      | ORIGIN                                                |
      | http://username@bitbucket.org/Originate/git-town.git  |
      | http://username@bitbucket.org/Originate/git-town      |
      | https://username@bitbucket.org/Originate/git-town.git |
      | https://username@bitbucket.org/Originate/git-town     |
      | git@bitbucket.org/Originate/git-town.git              |
      | git@bitbucket.org/Originate/git-town                  |


	Scenario Outline: origin includes path that looks like a URL
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is <ORIGIN>
    And I have "open" installed
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
		Then I see a new pull request with this url in my browser:
		  """
			https://bitbucket.org/Originate/originate.github.com/pull-request/new?dest=Originate%2Foriginate.github.com%3A%3Amain&source=Originate%2Foriginate.github.com%.*%3Afeature
			"""

		Examples:
      | ORIGIN                                                            |
      | http://username@bitbucket.org/Originate/originate.github.com.git  |
      | http://username@bitbucket.org/Originate/originate.github.com      |
      | https://username@bitbucket.org/Originate/originate.github.com.git |
      | https://username@bitbucket.org/Originate/originate.github.com     |
      | git@bitbucket.org/Originate/originate.github.com.git              |
      | git@bitbucket.org/Originate/originate.github.com                  |


	Scenario Outline: SSH style origin
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is <ORIGIN>
    And I have "open" installed
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
		Then I see a new pull request with this url in my browser:
		  """
			https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature
			"""

		Examples:
      | ORIGIN                                         |
      | ssh://git@bitbucket.org/Originate/git-town.git |
      | ssh://git@bitbucket.org/Originate/git-town     |
