Feature: git-town new-pull-request: when origin is an ssh identity

  When using my Git server via an SSH identity
  I want to be able to configure git town with the code hosting driver and origin hostname it should use
  So that new-pull-request works with my ssh identity


  Scenario Outline: ssh identity
    Given I have "open" installed
    And my repository has a feature branch named "feature"
    And my repo's remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "<ORIGIN_HOSTNAME>"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      <PULL_REQUEST_URL>
      """

    Examples:
      | ORIGIN_HOSTNAME | PULL_REQUEST_URL                                                                                                                        |
      | bitbucket.org   | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature  |
      | github.com      | https://github.com/Originate/git-town/compare/feature?expand=1                                                                          |
      | gitlab.com      | https://gitlab.com/Originate/git-town/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |
