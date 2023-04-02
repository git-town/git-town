Feature: display debug statistics

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |

  Scenario: result
    When I run "git-town rename-branch new --debug"
    Then it prints:
      """
      Ran 30 shell commands.
      """
    And the current branch is now "new"

  @this
  Scenario: undo
    Given I run "git-town rename-branch new"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 18 shell commands.
      """
    And the current branch is now "old"
