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
    And the current branch is "existing"
    When I run "git-town append new --beam" and enter into the dialog:
      | DIALOG                 | KEYS                             |
      | select commits 2 and 4 | down space down down space enter |

  @debug
  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                             |
      | existing | git fetch --prune --tags                                                                            |
      |          | git checkout main                                                                                   |
      | main     | git rebase origin/main --no-update-refs                                                             |
      |          | git checkout -b new                                                                                 |
      | new      | git cherry-pick {{ sha-before-run 'commit 2' }}                                                     |
      |          | git cherry-pick {{ sha-before-run 'commit 4' }}                                                     |
      |          | git checkout existing                                                                               |
      | existing | git rebase --onto {{ sha-before-run 'commit 4' }}^ {{ sha-before-run 'commit 4' }} --no-update-refs |
      |          | git rebase --onto {{ sha-before-run 'commit 2' }}^ {{ sha-before-run 'commit 2' }} --no-update-refs |
      |          | git push --force-with-lease --force-if-includes                                                     |
      |          | git checkout new                                                                                    |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | existing | local, origin | commit 1    |
      |          |               | commit 3    |
      | new      | local         | commit 2    |
      |          |               | commit 4    |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                |
      | new      | git checkout existing                                                  |
      | existing | git reset --hard {{ sha-before-run 'commit 4' }}                       |
      |          | git push --force-with-lease origin {{ sha 'initial commit' }}:existing |
      |          | git checkout main                                                      |
      | main     | git reset --hard {{ sha 'initial commit' }}                            |
      |          | git checkout existing                                                  |
      | existing | git branch -D new                                                      |
    And the initial commits exist now
    And the initial branches and lineage exist now
