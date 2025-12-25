@messyoutput
Feature: beam commits and uncommitted changes from a local branch onto a new feature branch and propose

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
    And the current branch is "existing"
    And tool "open" is installed
    And an uncommitted file
    And I ran "git add ."
    When I run "git-town hack new --beam --commit --message uncommitted --propose" and enter into the dialog:
      | DIALOG          | KEYS                             |
      | commits to beam | space down down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                 |
      | existing | git checkout -b new main                                                                                |
      | new      | git commit -m uncommitted                                                                               |
      |          | git cherry-pick {{ sha-initial 'commit 1' }}                                                            |
      |          | git cherry-pick {{ sha-initial 'commit 4' }}                                                            |
      |          | git checkout existing                                                                                   |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 4' }}^ {{ sha-initial 'commit 4' }} |
      |          | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1' }}^ {{ sha-initial 'commit 1' }} |
      |          | git checkout new                                                                                        |
      | new      | git push -u origin new                                                                                  |
      |          | open https://github.com/git-town/git-town/compare/new?expand=1&title=uncommitted                        |
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
      | existing | local         | commit 2    |
      |          |               | commit 3    |
      | new      | local, origin | uncommitted |
      |          |               | commit 1    |
      |          |               | commit 4    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                       |
      | existing | git reset --hard {{ sha-initial 'commit 4' }} |
      |          | git branch -D new                             |
      |          | git push origin :new                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
