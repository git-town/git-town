Feature: skip and park a perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE      | LOCATIONS     |
      | qa   | perennial | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE       | FILE NAME        | FILE CONTENT   |
      | qa     | local    | local commit  | conflicting_file | local content  |
      |        | origin   | origin commit | conflicting_file | origin content |
    And the current branch is "main"
    And I run "git-town sync --all"
    And Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git fetch --prune --tags                        |
      |        | git checkout qa                                 |
      | qa     | git -c rebase.updateRefs=false rebase origin/qa |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I run "git-town skip --park"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot park perennial branches
      """
    And a rebase is still in progress
    And branch "qa" still has type "perennial"
