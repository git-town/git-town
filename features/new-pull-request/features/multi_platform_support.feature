@skipWindows
Feature: support many browsers and operating systems

  Scenario Outline:
    Given the current branch is a feature branch "feature"
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And the "<TOOL>" tool is installed
    When I run "git-town new-pull-request"
    Then "<TOOL>" launches a new pull request with this url in my browser:
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
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has no tool to open browsers installed
    When I run "git-town new-pull-request"
    Then it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/compare/feature?expand=1
      """
