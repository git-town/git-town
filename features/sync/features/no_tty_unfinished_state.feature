Feature: descriptive error when there is unfinished state and no TTY is available

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: sync with unfinished state and no TTY
    When I run "git-town sync"
    Then Git Town prints the error:
      """
      there is an unfinished "sync" command and no interactive terminal available to resolve it.
      To continue the command, run: git town continue
      To skip the current branch, run: git town skip
      To undo the command, run: git town undo
      To discard the unfinished state, run: git town status reset
      """
