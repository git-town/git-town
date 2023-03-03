Feature: Describe a merge conflict during "git-town sync"

  @this
  Scenario: Git Town command in progress
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I run "git-town sync"
    When I run "git-town status"
    Then it prints something like:
      """
      The last Git Town command \(sync\) hit a problem .* ago.
      You can run "git town abort" to abort it.
      You can run "git town continue" to finish it.
      You can run "git town skip" to skip the currently failing step.
      """
    And it does not print "git town undo"
