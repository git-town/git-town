Feature: display debug statistics

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"

  Scenario: result
    When I run "git-town hack new --debug"
    Then it prints:
      """
      Ran 23 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town hack new"
    When I run "git town undo --debug"
    Then it prints:
      """
      Ran 13 shell commands.
      """
    And the current branch is now "main"
