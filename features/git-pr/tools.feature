Feature: git-pr with different tools

  Scenario Outline: supported tool installed
    Given I have a feature branch named "feature"
    And my remote origin is on GitHub through HTTPS ending with .git
    And I have <TOOL> installed
    And I am on the "feature" branch
    When I run `git pr`
    Then it opens a browser to a new GitHub pull request for the "feature" branch (<TOOL>)

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |


  Scenario: no supported tool installed
    Given I have a feature branch named "feature"
    And my remote origin is on GitHub through HTTPS ending with .git
    And I have nothing installed
    And I am on the "feature" branch
    When I run `git pr` while allowing errors
    Then I get the error "Opening a browser requires 'open' on Mac or 'xdg-open' on Linux."
