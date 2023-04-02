Feature: display debug statistics

  Background:
    Given the feature branches "active" and "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
      | old    | local, origin | old commit    |
    And origin deletes the "old" branch
    And the current branch is "old"

  Scenario: result
    When I run "git-town prune-branches --debug"
    Then it prints:
      """
      Ran 23 shell commands.
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | active | main   |

  Scenario: undo
    Given I run "git-town prune-branches"
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 11 shell commands.
      """
    And the current branch is now "old"
    And the initial branches and hierarchy exist
