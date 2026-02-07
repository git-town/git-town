Feature: provide a diff-filter

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     | FILE NAME         | FILE CONTENT |
      | main   | local    | main commit | existing-file.txt | main content |
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS |
      | feature-1 | feature | main   | local     |
      | feature-2 | feature | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE  | FILE NAME         | FILE CONTENT      |
      | feature-1 | local    | commit 1 | existing-file.txt | feature-1 content |
      | feature-1 | local    | commit 1 | new-file.txt      | new content       |
      | feature-2 | local    | commit 2 | other-file.txt    | content           |
    And the current branch is "feature-1"
    When I run "git-town diff-parent --diff-filter=A"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                              |
      | feature-1 | git diff --diff-filter=A --merge-base main feature-1 |
    And Git Town prints:
      """
      new-file
      """
    And Git Town does not print "existing-file.txt"
