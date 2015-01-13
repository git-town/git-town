Feature: git-pr: multi-platform support

  Scenario Outline: supported tool installed
    Given I have a feature branch named "feature"
    And my remote origin is https://github.com/Originate/git-town.git
    And I have "<TOOL>" installed
    And I am on the "feature" branch
    When I run `git pr`
    Then I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |


  Scenario: no supported tool installed
    Given I have a feature branch named "feature"
    And my remote origin is https://github.com/Originate/git-town.git
    And I have no command that opens browsers installed
    And I am on the "feature" branch
    When I run `git pr` while allowing errors
    Then I get the error "Opening a browser requires 'open' on Mac or 'xdg-open' on Linux."
