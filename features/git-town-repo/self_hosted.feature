Feature: git-town repo: when origin is a self hosted servie

  Scenario Outline: self hosted
    Given I have "open" installed
    And my repo's remote origin is "git@self-hosted:Originate/git-town.git"
    And I configure "git-town.code-hosting-driver" as "<DRIVER>"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      <REPO_URL>
      """

    Examples:
      | DRIVER    | REPO_URL                               |
      | bitbucket | https://self-hosted/Originate/git-town |
      | github    | https://self-hosted/Originate/git-town |
      | gitlab    | https://self-hosted/Originate/git-town |
