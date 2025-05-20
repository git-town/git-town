Feature: undo changes made manually

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --all"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-2 | git checkout branch-1 |
    And Git Town prints:
      """
      Run "git town continue" to go to the next branch.
      """
    And I add this commit to the current branch:
      | MESSAGE  | FILE NAME | FILE CONTENT |
      | commit 1 | file_1    | content 1    |
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-1 | git checkout branch-2 |
    And Git Town prints:
      """
      Run "git town continue" to go to the next branch.
      """
    And the current branch is now "branch-2"
    And I add this commit to the current branch:
      | MESSAGE  | FILE NAME | FILE CONTENT |
      | commit 2 | file_2    | content 2    |
    When I run "git-town continue"
    And Git Town prints:
      """
      Branch walk done.
      """

  Scenario: result
    Then these commits exist now
      | BRANCH   | LOCATION | MESSAGE  |
      | branch-1 | local    | commit 1 |
      | branch-2 | local    | commit 2 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                     |
      | branch-2 | git checkout branch-1                       |
      | branch-1 | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout branch-2                       |
      | branch-2 | git reset --hard {{ sha 'initial commit' }} |
    And the current branch is now "branch-2"
    And no commits exist now
