Feature: a grandchild branch has conflicts while its parent was deleted remotely

  Background:
    Given the current branch is a feature branch "child"
    And a feature branch "grandchild" as a child of "child"
    And the commits
      | BRANCH     | LOCATION | MESSAGE                       | FILE NAME        | FILE CONTENT       |
      | main       | local    | conflicting main commit       | conflicting_file | main content       |
      | grandchild | local    | conflicting grandchild commit | conflicting_file | grandchild content |
    And origin deletes the "child" branch
    When I run "git-town sync --all"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND   |
      | child  | git fetch |
    And it prints the error:
      """
      exit status 1
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git-town continue".
      To go back to where you started, run "git-town undo".
      To continue by skipping the current branch, run "git-town skip".
      """
    And the current branch is now "feature"
    And the uncommitted file is stashed
    And a merge is now in progress
