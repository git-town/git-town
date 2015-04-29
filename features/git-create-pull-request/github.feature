Feature: git-create-pull-request when origin is on GitHub

  As a developer having finished a feature in a repository hosted on GitHub
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Scenario Outline: result
    Given I have a feature branch named "feature"
    And my remote origin is <ORIGIN>
    And I have "open" installed
    And I am on the "feature" branch
    When I run `git create-pull-request`
    Then I see a new GitHub pull request for the "feature" branch in the "<REPOSITORY>" repo in my browser

    Examples:
      | ORIGIN                                                | REPOSITORY                     |
      | http://github.com/Originate/git-town.git              | Originate/git-town             |
      | http://github.com/Originate/git-town                  | Originate/git-town             |
      | https://github.com/Originate/git-town.git             | Originate/git-town             |
      | https://github.com/Originate/git-town                 | Originate/git-town             |
      | git@github.com:Originate/git-town.git                 | Originate/git-town             |
      | git@github.com:Originate/git-town                     | Originate/git-town             |
      | http://github.com/Originate/originate.github.com.git  | Originate/originate.github.com |
      | http://github.com/Originate/originate.github.com      | Originate/originate.github.com |
      | https://github.com/Originate/originate.github.com.git | Originate/originate.github.com |
      | https://github.com/Originate/originate.github.com     | Originate/originate.github.com |
      | git@github.com:Originate/originate.github.com.git     | Originate/originate.github.com |
      | git@github.com:Originate/originate.github.com         | Originate/originate.github.com |
