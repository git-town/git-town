Feature: displaying the branches in the middle of an ongoing sync merge conflict

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    And Git Town setting "sync-feature-strategy" is "compress"
    And an uncommitted file
    And I ran "git-town sync"
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I run "git-town branch"

  @this
  Scenario: result
    Then it runs no commands
    And it prints:
      """
        main
          feature
      """
