Feature: undo a hack after additional commits to the new branch

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"
    When I run "git-town hack new"
    And I add commit "additional commit" to the "new" branch

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot undo because these branches received new commits after the last Git Town command: new
      Undoing would delete these commits. To undo anyway, run "git town undo --force".
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE           |
      | main   | local, origin | main commit       |
      | new    | local         | additional commit |

  Scenario: undo with force
    When I run "git-town undo --force"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | new    | git checkout main                           |
      | main   | git reset --hard {{ sha 'initial commit' }} |
      |        | git branch -D new                           |
    And Git Town prints:
      """
      deleted branch new
      """
    And the initial commits exist now
    And the initial branches and lineage exist now
