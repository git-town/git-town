Feature: get only the names of changed files

  Background:
    Given a Git repo with origin

  Scenario: feature branch
    Given the branches
      | NAME      | TYPE    | PARENT | LOCATIONS |
      | feature-1 | feature | main   | local     |
      | feature-2 | feature | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE   | FILE NAME   |
      | feature-1 | local    | commit 1A | file-1A.txt |
      | feature-1 | local    | commit 1B | file-1B.txt |
      | feature-2 | local    | commit 2  | file-2.txt  |
    And the current branch is "feature-1"
    When I run "git-town diff-parent --name-only"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                          |
      | feature-1 | git diff --name-only --merge-base main feature-1 |
    And Git Town prints:
      """
      file-1A.txt
      file-1B.txt
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
