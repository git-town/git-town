@skipWindows
Feature: handle conflicts between the current feature branch and the main branch (in a local repo)

  Background:
    Given a local Git repo
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And I ran "git-town hack feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |         | git checkout --theirs conflicting_file                                       |
      |         | git add conflicting_file                                                     |
      |         | GIT_EDITOR=true git rebase --continue                                        |
    And no rebase is now in progress
    And all branches are now synchronized
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                 |
      | feature | git reset --hard {{ sha 'conflicting feature commit' }} |
    And no rebase is now in progress
    And the initial commits exist now
