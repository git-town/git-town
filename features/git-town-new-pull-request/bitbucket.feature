Feature: git-new-pull-request when origin is on Bitbucket

  As a developer having finished a feature in a repository hosted on Bitbucket
  I want to be able to quickly create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Scenario Outline: creating pull-requests
    Given I have a feature branch named "feature"
    And my remote origin is <ORIGIN>
    And I have "open" installed
    And I am on the "feature" branch
    When I run `git town-new-pull-request`
    Then I see a new Bitbucket pull request for the "feature" branch in the "<REPOSITORY>" repo in my browser

    Examples:
      | ORIGIN                                                            | REPOSITORY                     |
      | http://username@bitbucket.org/Originate/git-town.git              | Originate/git-town             |
      | http://username@bitbucket.org/Originate/git-town                  | Originate/git-town             |
      | https://username@bitbucket.org/Originate/git-town.git             | Originate/git-town             |
      | https://username@bitbucket.org/Originate/git-town                 | Originate/git-town             |
      | git@bitbucket.org/Originate/git-town.git                          | Originate/git-town             |
      | git@bitbucket.org/Originate/git-town                              | Originate/git-town             |
      | http://username@bitbucket.org/Originate/originate.github.com.git  | Originate/originate.github.com |
      | http://username@bitbucket.org/Originate/originate.github.com      | Originate/originate.github.com |
      | https://username@bitbucket.org/Originate/originate.github.com.git | Originate/originate.github.com |
      | https://username@bitbucket.org/Originate/originate.github.com     | Originate/originate.github.com |
      | git@bitbucket.org/Originate/originate.github.com.git              | Originate/originate.github.com |
      | git@bitbucket.org/Originate/originate.github.com                  | Originate/originate.github.com |
