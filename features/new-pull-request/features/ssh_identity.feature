@skipWindows
Feature: use a SSH identity

  Scenario Outline: ssh identity
    And my computer has the "open" tool installed
    And a feature branch "feature"
    And my repo's origin is "git@my-ssh-identity:git-town/git-town.git"
    And the "code-hosting-origin-hostname" setting is "<ORIGIN_HOSTNAME>"
    And the current branch is "feature"
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
