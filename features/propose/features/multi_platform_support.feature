@skipWindows
Feature: support many browsers and operating systems

  Scenario Outline:
    Given the current branch is a feature branch "feature"
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "<TOOL>" is installed
    When I run "git-town propose"
    Then "<TOOL>" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """

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
    Given the current branch is a feature branch "feature"
    And the origin is "https://github.com/git-town/git-town.git"
    And no tool to open browsers is installed
    When I run "git-town propose"
    Then it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/compare/feature?expand=1
      """
