Feature: multi-platform support

  @skipWindows
  Scenario Outline: supported tool installed
    Given my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has the "<TOOL>" tool installed
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town
      """

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |

  @skipWindows
  Scenario: no supported tool installed
    Given my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has no tool to open browsers installed
    When I run "git-town repo"
    Then it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town
      """
