Feature: display debug and performance statistics

  Background:
    Given the feature branches "active" and "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
      | old    | local, origin | old commit    |
    And origin deletes the "old" branch
    And the current branch is "old"
    When I run "git-town prune-branches --debug"

  Scenario: result
    Then it prints:
      """
      Ran 20 shell commands.
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | active | main   |

  @this
  Scenario: undo
    When I run "git-town undo --debug"
    Then it prints:
      """
      Ran 9 shell commands.
      """
    And the current branch is now "old"
    And the initial branches and hierarchy exist
