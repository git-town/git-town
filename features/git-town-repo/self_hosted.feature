Feature: git-town repo: when origin is a self hosted servie

  Background:
    Given I have "open" installed


  Scenario: self hosted bitbucket
    Given my remote origin is "git@self-hosted-bitbucket:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "bitbucket"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      https://self-hosted-bitbucket/Originate/git-town
      """


  Scenario: self hosted github
    Given my remote origin is "git@self-hosted-github:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "github"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      https://self-hosted-github/Originate/git-town
      """


  Scenario: self hosted gitlab
    Given my remote origin is "git@self-hosted-gitlab:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "gitlab"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      https://self-hosted-gitlab/Originate/git-town
      """
