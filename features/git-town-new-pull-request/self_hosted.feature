Feature: git-town new-pull-request: when origin is a self hosted servie

  When creating a pull request on a self-hosted repository whose type cannot be derived from its URL
  I want to specify the type of the hosted repository
  So that Git Town can support specific features of my repository platform.


  Scenario Outline: self hosted
    Given I have "open" installed
    And my repository has a feature branch named "feature"
    And my repo's remote origin is "git@self-hosted:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "<DRIVER>"
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with this url in my browser:
      """
      <PULL_REQUEST_URL>
      """

    Examples:
      | DRIVER    | PULL_REQUEST_URL                                                                                                                         |
      | bitbucket | https://self-hosted/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature     |
      | github    | https://self-hosted/Originate/git-town/compare/feature?expand=1                                                                          |
      | gitlab    | https://self-hosted/Originate/git-town/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |
