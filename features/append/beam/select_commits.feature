@messyoutput
Feature: beam multiple commits onto a new child branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE     |
      | main     | origin   | main commit |
      | existing | local    | commit 1    |
      | existing | local    | commit 2    |
      | existing | local    | commit 3    |
      | existing | local    | commit 4    |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "existing"
    When I run "git-town append new --beam" and enter into the dialog:
      | DIALOG          | KEYS                             |
      | commits to beam | space down down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                 |
      | existing | git checkout -b new                                                                                     |
      | new      | git checkout existing                                                                                   |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 4' }}^ {{ sha-initial 'commit 4' }} |
      |          | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1' }}^ {{ sha-initial 'commit 1' }} |
      |          | git push --force-with-lease --force-if-includes                                                         |
      |          | git checkout new                                                                                        |
      | new      | git -c rebase.updateRefs=false rebase existing                                                          |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        existing
          new
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local, origin | commit 2    |
      |          |               | commit 3    |
      | new      | local         | commit 1    |
      |          |               | commit 4    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                |
      | new      | git checkout existing                                                  |
      | existing | git reset --hard {{ sha-initial 'commit 4' }}                          |
      |          | git push --force-with-lease origin {{ sha 'initial commit' }}:existing |
      |          | git branch -D new                                                      |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: amend the beamed commit
    And I amend this commit
      | BRANCH | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT    |
      | new    | local    | commit 4b | file_4    | amended content |
    And the current branch is "new"
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                         |
      | new      | git fetch --prune --tags                                                                        |
      |          | git checkout main                                                                               |
      | main     | git -c rebase.updateRefs=false rebase origin/main                                               |
      |          | git checkout existing                                                                           |
      | existing | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }}                    |
      |          | git push --force-with-lease --force-if-includes                                                 |
      |          | git checkout new                                                                                |
      | new      | git -c rebase.updateRefs=false rebase --onto existing {{ sha-in-origin-before-run 'commit 3' }} |
      |          | git push -u origin new                                                                          |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | existing | local, origin | commit 2    |
      |          |               | commit 3    |
      | new      | local, origin | commit 1    |
      |          |               | commit 4b   |
