Feature: get only the names of changed files

  Background:
    Given a Git repo with origin

  @this
  Scenario: feature branch
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE              | FILE NAME |
      | feature | local    | local feature commit | file.txt  |
    And the current branch is "feature"
    When I run "git-town diff-parent --name-only"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git diff --merge-base main feature |
    And Git Town prints:
      """
      file.txt
      """

  Scenario: child branch
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | parent | feature | main   | local     |
      | child  | feature | parent | local     |
    And the current branch is "child"
    When I run "git-town diff-parent"
    Then Git Town runs the commands
      | BRANCH | COMMAND                            |
      | child  | git diff --merge-base parent child |
