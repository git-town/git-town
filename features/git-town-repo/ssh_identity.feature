Feature: git-town repo: when origin is an ssh identity

  Scenario Outline: ssh identity
    Given I have "open" installed
    And my repo's remote origin is "git@my-ssh-identity:Originate/git-town.git"
    And I configure "git-town.code-hosting-origin-hostname" as "<ORIGIN_HOSTNAME>"
    When I run `git-town repo`
    Then I see my repo homepage this url in my browser:
      """
      <REPO_URL>
      """

    Examples:
      | ORIGIN_HOSTNAME | REPO_URL                                 |
      | bitbucket.org   | https://bitbucket.org/Originate/git-town |
      | github.com      | https://github.com/Originate/git-town    |
      | gitlab.com      | https://gitlab.com/Originate/git-town    |
