Feature: undo changes made manually

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
      | branch-2 | feature | main   | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --all -- git commit --allow-empty -m commit"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                            |
      | branch-2 | git checkout branch-1              |
      | branch-1 | git commit --allow-empty -m commit |
      |          | git checkout branch-2              |
      | branch-2 | git commit --allow-empty -m commit |
    And Git Town prints:
      """
      Branch walk done.
      """

  Scenario: result
    Then these commits exist now
      | BRANCH   | LOCATION | MESSAGE |
      | branch-1 | local    | commit  |
      | branch-2 | local    | commit  |

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
