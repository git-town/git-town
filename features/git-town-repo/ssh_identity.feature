Feature: git-town repo: when origin is an ssh identity

  Background:
    Given I have "open" installed


  Scenario: bitbucket ssh identity
    Given my remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "bitbucket.org"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      https://bitbucket.org/Originate/git-town
      """


  Scenario: github ssh identity
    Given my remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "github.com"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      https://github.com/Originate/git-town
      """


  Scenario: gitlab ssh identity
    Given my remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "gitlab.com"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      https://gitlab.com/Originate/git-town
      """
