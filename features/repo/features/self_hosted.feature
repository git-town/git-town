@skipWindows
Feature: self hosted servie

  Scenario Outline:
    Given my computer has the "open" tool installed
    And my repo's origin is "git@self-hosted:git-town/git-town.git"
    And my repo has "git-town.code-hosting-driver" set to "<DRIVER>"
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      <REPO_URL>
      """

    Examples:
      | DRIVER    | REPO_URL                              |
      | bitbucket | https://self-hosted/git-town/git-town |
      | github    | https://self-hosted/git-town/git-town |
      | gitea     | https://self-hosted/git-town/git-town |
      | gitlab    | https://self-hosted/git-town/git-town |
