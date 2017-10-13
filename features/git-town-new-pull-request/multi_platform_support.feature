Feature: git-new-pull-request: multi-platform support

  Scenario Outline: supported tool installed
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is https://github.com/Originate/git-town.git
    And I have "<TOOL>" installed
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser

    Examples:
      | TOOL          |
      | open          |
      | xdg-open      |
      | cygstart      |
      | x-www-browser |
      | firefox       |
      | opera         |
      | mozilla       |
      | netscape      |


  Scenario: no supported tool installed
    Given my repository has a feature branch named "feature"
    And my repo's remote origin is https://github.com/Originate/git-town.git
    And I have no command that opens browsers installed
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then Git Town prints the error "Cannot open a browser"
