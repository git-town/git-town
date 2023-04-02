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
      Ran 30 shell commands.
      """
    And the current branch is now "new"
    And now these commits exist
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
      | new      | local         | existing commit |
    And this branch hierarchy exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    Given I run "git-town append new"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 15 shell commands.
      """
    And the current branch is now "existing"
    And now the initial commits exist
    And the initial branch hierarchy exists
