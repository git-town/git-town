Feature: git-town new-pull-request: when origin is an ssh identity

  As a developer trying to create a pull request in a that uses an SSH identity
  I want to be able to configure git town with the code hosting driver and origin hostname it should use
  So that new-pull-request works with my ssh identity


  Background:
    Given I have "open" installed


  Scenario: bitbucket ssh identity
    Given I have a feature branch named "feature"
    And my remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "bitbucket.org"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature
      """


  Scenario: github ssh identity
    Given I have a feature branch named "feature"
    And my remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "github.com"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      https://github.com/Originate/git-town/compare/feature?expand=1
      """


  Scenario: gitlab ssh identity
    Given I have a feature branch named "feature"
    And my remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "gitlab.com"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      https://gitlab.com/Originate/git-town/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main
      """
