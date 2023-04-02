Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |

  Scenario: result
    When I run "git-town ship -m done --debug"
    Then it prints:
      """
      Ran 46 shell commands.
      """
    And the current branch is now "main"

  Scenario: undo
    Given I run "git-town ship -m done"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 19 shell commands.
      """
    And the current branch is now "feature"
