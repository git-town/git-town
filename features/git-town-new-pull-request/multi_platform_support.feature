Feature: git-new-pull-request: multi-platform support

  Scenario Outline: supported tool installed
    Given my repository has a feature branch named "feature"
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has the "<TOOL>" tool installed
    And I am on the "feature" branch
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
    Given my repository has a feature branch named "feature"
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has no tool to open browsers installed
    And I am on the "feature" branch
    When I run "git-town new-pull-request"
    Then it prints the error:
      """
      Cannot open a browser
      """
