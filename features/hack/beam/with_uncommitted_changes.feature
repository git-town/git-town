@messyoutput
Feature: beam a commit and uncommitted changes onto a new feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local, origin | commit 1    |
      | existing | local, origin | commit 2    |
      | existing | local, origin | commit 3    |
    And the current branch is "existing"
    And an uncommitted file
    And I ran "git add ."
    When I run "git-town hack new --beam --commit --message uncommitted" and enter into the dialog:
      | DIALOG          | KEYS             |
      | commits to beam | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                 |
      | existing | git checkout -b new main                                                                                |
      | new      | git commit -m uncommitted                                                                               |
      |          | git cherry-pick {{ sha-initial 'commit 2' }}                                                            |
      |          | git checkout existing                                                                                   |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 2' }}^ {{ sha-initial 'commit 2' }} |
      |          | git push --force-with-lease --force-if-includes                                                         |
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
      | existing | local, origin | commit 1    |
      |          |               | commit 3    |
      | new      | local         | uncommitted |
      |          |               | commit 2    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | existing | git reset --hard {{ sha-initial 'commit 3' }}   |
      |          | git push --force-with-lease --force-if-includes |
      |          | git branch -D new                               |
    And the initial branches and lineage exist now
    And the initial commits exist now
