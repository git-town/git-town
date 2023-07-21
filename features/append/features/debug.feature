Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |

  Scenario: result
    When I run "git-town append new --debug"
    Then it prints:
      """
      Ran 29 shell commands.
      """
    And the current branch is now "new"

  Scenario: undo
    Given I ran "git-town append new"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 14 shell commands.
      """
    And the current branch is now "existing"
