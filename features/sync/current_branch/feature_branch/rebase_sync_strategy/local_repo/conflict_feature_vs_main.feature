@skipWindows
Feature: handle conflicts between the current feature branch and the main branch (in a local repo)

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    And a rebase is now in progress
    And file "conflicting_file" now has content:
      """
      <<<<<<< HEAD
      main content
      =======
      feature content
      >>>>>>> {{ sha-short 'conflicting feature commit' }} (conflicting feature commit)
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND            |
      | feature | git rebase --abort |
    And no rebase is now in progress
    And the initial commits exist now

  Scenario: continue without resolving the conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file" with "main and feature content"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                               |
      | feature | GIT_EDITOR=true git rebase --continue |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT             |
      | main    | local    | conflicting main commit    | conflicting_file | main content             |
      | feature | local    | conflicting feature commit | conflicting_file | main and feature content |

  Scenario: resolve, rebase, and continue
    When I resolve the conflict in "conflicting_file" with "main and feature content"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then Git Town runs no commands
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT             |
      | main    | local    | conflicting main commit    | conflicting_file | main content             |
      | feature | local    | conflicting feature commit | conflicting_file | main and feature content |
