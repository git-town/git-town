Feature: git-repo with different tools

  Scenario Outline: with a supported tool
    Given my remote origin is on GitHub through HTTPS ending with .git
    And I have <TOOL> installed
    When I run `git repo`
    Then it opens a browser to the homepage of my GitHub repository (<TOOL>)

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |


  Scenario: without a supported tool
    Given my remote origin is on GitHub through HTTPS ending with .git
    And I have nothing installed
    When I run `git repo` while allowing errors
    Then I get the error "Opening a browser requires 'open' on Mac or 'xdg-open' on Linux."
