Feature: commit down

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT   | LOCATIONS     |
      | branch-1 | observed | main     | local, origin |
      | branch-2 | feature  | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | branch-1 | local, origin | commit 1a | file_1    | content 1    |
      | branch-2 | local, origin | commit 2a | file_2    | content 2    |
    And the current branch is "branch-2"
    And an uncommitted file "changes" with content "my changes"
    And I ran "git add changes"
    When I run "git-town commit --down -m commit-1b"

  @this
  Scenario: result
    # Then Git Town runs the commands
    #   | BRANCH | COMMAND |
    And Git Town prints the error:
      """
      cannot commit into branch branch-1 because it is an observed branch
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
