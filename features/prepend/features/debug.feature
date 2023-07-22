Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |

  Scenario: result
    When I run "git-town prepend parent --debug"
    Then it prints:
      """
      Ran 26 shell commands.
      """
    And the current branch is now "parent"

  Scenario: undo
    Given I ran "git-town prepend parent"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 15 shell commands.
      """
    And the current branch is now "old"
