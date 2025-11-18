@messyoutput
Feature: beam commits from a branch on a worktree different from main

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE   | PARENT | LOCATIONS     |
      | existing-1 | (none) |        | local, origin |
    And the commits
      | BRANCH     | LOCATION | MESSAGE     |
      | main       | origin   | main commit |
      | existing-1 | local    | commit 1a   |
      | existing-1 | local    | commit 1b   |
    And the current branch is "existing-1"
    And branch "main" is active in another worktree
    When I run "git-town hack new --beam" and enter into the dialog:
      | DIALOG                         | KEYS             |
      | parent branch for "existing-1" | enter            |
      | commits to beam                | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                                                                   |
      | existing-1 | git checkout -b new main                                                                                  |
      | new        | git cherry-pick {{ sha-initial 'commit 1b' }}                                                             |
      |            | git checkout existing-1                                                                                   |
      | existing-1 | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1b' }}^ {{ sha-initial 'commit 1b' }} |
      |            | git push --force-with-lease --force-if-includes                                                           |
      |            | git checkout new                                                                                          |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        existing-1
        new
      """
    And these commits exist now
      | BRANCH     | LOCATION                | MESSAGE     |
      | main       | origin                  | main commit |
      | existing-1 | local, origin, worktree | commit 1a   |
      | new        | local                   | commit 1b   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                                  |
      | new        | git checkout existing-1                                                  |
      | existing-1 | git reset --hard {{ sha-initial 'commit 1b' }}                           |
      |            | git push --force-with-lease origin {{ sha 'initial commit' }}:existing-1 |
      |            | git branch -D new                                                        |
    And the initial branches and lineage exist now
    And the initial commits exist now
