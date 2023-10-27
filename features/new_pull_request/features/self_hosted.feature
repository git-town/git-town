@skipWindows
Feature: self-hosted service

  Scenario Outline: self hosted
    Given tool "open" is installed
    And the origin is "git@self-hosted:git-town/git-town.git"
    And Git Town setting "code-hosting-driver" is "<DRIVER>"
    And the current branch is a feature branch "feature"
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      <PULL_REQUEST_URL>
      """

    Examples:
      | DRIVER    | PULL_REQUEST_URL                                                                                                                          |
      | bitbucket | https://self-hosted/git-town/git-town/pull-request/new?dest=git-town%2Fgit-town%3A%3Amain&source=git-town%2Fgit-town%.*%3Afeature         |
      | github    | https://self-hosted/git-town/git-town/compare/feature?expand=1                                                                            |
      | gitea     | https://self-hosted/git-town/git-town/compare/main...feature                                                                              |
      | gitlab    | https://self-hosted/git-town/git-town/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |
