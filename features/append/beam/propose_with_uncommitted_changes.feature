@messyoutput
Feature: beam a commit and uncommitted changes onto a new child branch and propose

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | origin        | main commit |
      | existing | local, origin | commit 1    |
      | existing | local, origin | commit 2    |
      | existing | local, origin | commit 3    |
      | existing | local, origin | commit 4    |
      | existing | local, origin | commit 5    |
    And the current branch is "existing"
    And tool "open" is installed
    And an uncommitted file
    And I ran "git add ."
    When I run "git-town append new --beam --commit --message uncommitted --propose" and enter into the dialog:
      | DIALOG          | KEYS                             |
      | commits to beam | space down down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                 |
      | existing | git checkout -b new                                                                                     |
      | new      | git commit -m uncommitted                                                                               |
      |          | git checkout existing                                                                                   |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 4' }}^ {{ sha-initial 'commit 4' }} |
      |          | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1' }}^ {{ sha-initial 'commit 1' }} |
      |          | git push --force-with-lease --force-if-includes                                                         |
      |          | git checkout new                                                                                        |
      | new      | git -c rebase.updateRefs=false rebase existing                                                          |
      |          | git push -u origin new                                                                                  |
      |          | Finding proposal from new into existing ...                                                             |
      |          | open https://github.com/git-town/git-town/compare/existing...new?expand=1&title=uncommitted             |
      |          | git checkout existing                                                                                   |
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
      |          |               | commit 5    |
      | new      | local, origin | commit 1    |
      |          |               | commit 4    |
      |          |               | uncommitted |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | existing | git reset --hard {{ sha-initial 'commit 5' }}   |
      |          | git push --force-with-lease --force-if-includes |
      |          | git branch -D new                               |
      |          | git push origin :new                            |
    And the initial branches and lineage exist now
    And the initial commits exist now
