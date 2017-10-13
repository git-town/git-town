Feature: git-repo: multi-platform support

  Scenario Outline: supported tool installed
    Given my repo's remote origin is https://github.com/Originate/git-town.git
    And I have "<TOOL>" installed
    When I run `git-town repo`
    Then I see the GitHub homepage of the "Originate/git-town" repository in my browser

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |


  Scenario: no supported tool installed
    Given my repo's remote origin is https://github.com/Originate/git-town.git
    And I have no command that opens browsers installed
    When I run `git-town repo`
    Then Git Town prints the error "Cannot open a browser"
