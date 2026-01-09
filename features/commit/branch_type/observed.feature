Feature: commit down into an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT   | LOCATIONS     |
      | branch-1 | observed | main     | local, origin |
      | branch-2 | feature  | branch-1 | local, origin |
    And the current branch is "branch-2"
    And an uncommitted file "changes" with content "my changes"
    And I ran "git add changes"
    When I run "git-town commit --down -m commit-1b"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot commit into branch branch-1 because it is an observed branch
      """
  # Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
