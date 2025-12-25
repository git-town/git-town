Feature: displaying the repo in the middle of an ongoing sync merge conflict

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    And tool "open" is installed
    And I ran "git-town sync" and ignore the error
    When I run "git-town repo"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                   |
      | feature | open https://github.com/git-town/git-town |
