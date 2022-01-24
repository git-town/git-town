Feature: git-town new-pull-request: when origin is an ssh identity

  When using my Git server via an SSH identity
  I want to be able to configure git town with the code hosting driver and origin hostname it should use
  So that new-pull-request works with my ssh identity

  @skipWindows
  Scenario Outline: ssh identity
    And my computer has the "open" tool installed
    And my repo has a feature branch named "feature"
    And my repo's origin is "git@my-ssh-identity:git-town/git-town.git"
    And my repo has "git-town.code-hosting-origin-hostname" set to "<ORIGIN_HOSTNAME>"
    And I am on the "feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      <PULL_REQUEST_URL>
      """

    Examples:
      | ORIGIN_HOSTNAME | PULL_REQUEST_URL                                                                                                                       |
      | bitbucket.org   | https://bitbucket.org/git-town/git-town/pull-request/new?dest=git-town%2Fgit-town%3A%3Amain&source=git-town%2Fgit-town%.*%3Afeature    |
      | github.com      | https://github.com/git-town/git-town/compare/feature?expand=1                                                                          |
      | gitlab.com      | https://gitlab.com/git-town/git-town/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |
