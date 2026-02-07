Feature: provide the diff filter

  Background:
    Given a Git repo with origin

  @this
  Scenario: feature branch
    Given the branches
      | NAME      | TYPE    | PARENT | LOCATIONS |
      | feature-1 | feature | main   | local     |
      | feature-2 | feature | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE  | FILE NAME  |
      | feature-1 | local    | commit 1 | file-1.txt |
      | feature-2 | local    | commit 2 | file-2.txt |
    And the current branch is "feature-1"
    And an uncommitted file "new-file" with content "content"
    When I run "git-town diff-parent --diff-filter=A"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                              |
      | feature-1 | git diff --diff-filter=A --merge-base main feature-1 |
    And Git Town prints:
      """
      new-file
      """
    And Git Town does not print "file-2.txt"

  Scenario: child branch
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | parent | feature | main   | local     |
      | child  | feature | parent | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE   | FILE NAME   |
      | parent | local    | commit 1A | file-1.txt  |
      | child  | local    | commit 2A | file-2A.txt |
      | child  | local    | commit 2B | file-2B.txt |
    And the current branch is "child"
    When I run "git-town diff-parent --name-only"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                        |
      | child  | git diff --name-only --merge-base parent child |
    And Git Town prints:
      """
      file-2A.txt
      file-2B.txt
      """
    And Git Town does not print "file-1.txt"
