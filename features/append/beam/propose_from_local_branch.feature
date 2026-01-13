@messyoutput
Feature: beam commits and uncommitted changes from a local branch onto a new child branch and propose

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE     |
      | main     | origin   | main commit |
      | existing | local    | commit 1    |
      | existing | local    | commit 2    |
      | existing | local    | commit 3    |
      | existing | local    | commit 4    |
      | existing | local    | commit 5    |
    And the current branch is "existing"
    And tool "open" is installed
    And an uncommitted file
    And I ran "git add ."
    When I run "git-town append new --beam --commit --message uncommitted --propose" and enter into the dialog:
      | DIALOG          | KEYS                             |
      | commits to beam | space down down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                                         |
      | existing | git checkout -b new                                                                                                             |
      | new      | git commit -m uncommitted                                                                                                       |
      |          | git checkout existing                                                                                                           |
      | existing | git -c rebase.updateRefs=false rebase --onto 8a201b280df8bdaaf28b7a7a1062842531f873d5^ 8a201b280df8bdaaf28b7a7a1062842531f873d5 |
      |          | git -c rebase.updateRefs=false rebase --onto 1f3ca057cfbf5925f4755c0358f508410df05786^ 1f3ca057cfbf5925f4755c0358f508410df05786 |
      |          | git checkout new                                                                                                                |
      | new      | git -c rebase.updateRefs=false rebase existing                                                                                  |
      |          | git push -u origin new                                                                                                          |
      |          | Finding proposal from new into existing ...                                                                                     |
      |          | open https://github.com/git-town/git-town/compare/existing...new?expand=1&title=uncommitted                                     |
      |          | git checkout existing                                                                                                           |
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
      | existing | local         | commit 2    |
      |          |               | commit 3    |
      |          |               | commit 5    |
      | new      | local, origin | commit 1    |
      |          |               | commit 4    |
      |          |               | uncommitted |
      |          | origin        | commit 2    |
      |          |               | commit 3    |
      |          |               | commit 5    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                       |
      | existing | git reset --hard {{ sha-initial 'commit 5' }} |
      |          | git branch -D new                             |
      |          | git push origin :new                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
