Feature: git-repo with different tools

  Scenario Outline: supported tool installed
    Given my remote origin is on GitHub through HTTPS ending with .git
    And I have "<TOOL>" installed
    When I run `git repo`
    Then I see a browser open to the homepage of my GitHub repository (<TOOL>)

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |


  Scenario: no supported tool installedol
    Given my remote origin is on GitHub through HTTPS ending with .git
    And I have no command that opens browsers installed
    When I run `git repo` while allowing errors
    Then I get the error "Opening a browser requires 'open' on Mac or 'xdg-open' on Linux."
