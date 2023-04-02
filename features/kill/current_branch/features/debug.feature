Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "current"
    And a feature branch "other"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |

  Scenario: result
    When I run "git-town kill --debug"
    Then it prints:
      """
      Ran 29 shell commands.
      """
    And the current branch is now "main"

  @this
  Scenario: undo
    Given I ran "git-town kill"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 12 shell commands.
      """
    And the current branch is now "current"
