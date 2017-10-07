Feature: git-town new-pull-request: when origin is a self hosted servie

  As a developer trying to create a pull request when using a self hosted service
  I want to be able to configure git town with the code hosting driver it should use
  So that new-pull-request works with my self hosted service


  Background:
    Given I have "open" installed


  Scenario: self hosted bitbucket
    Given I have a feature branch named "feature"
    And my remote origin is "git@self-hosted-bitbucket:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "bitbucket"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      https://self-hosted-bitbucket/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature
      """


  Scenario: self hosted github
    Given I have a feature branch named "feature"
    And my remote origin is "git@self-hosted-github:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "github"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      https://self-hosted-github/Originate/git-town/compare/feature?expand=1
      """


  Scenario: self hosted gitlab
    Given I have a feature branch named "feature"
    And my remote origin is "git@self-hosted-gitlab:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "gitlab"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      https://self-hosted-gitlab/Originate/git-town/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main
      """
