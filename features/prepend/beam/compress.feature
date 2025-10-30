@messyoutput
Feature: prepend a branch to a feature branch using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | old    | local, origin | commit 1 | file 1    | content 1    |
      | old    | local, origin | commit 2 | file 2    | content 2    |
      | old    | local, origin | commit 3 | file 3    | content 3    |
      | old    | local, origin | commit 4 | file 4    | content 4    |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "old"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town prepend parent --beam" and enter into the dialog:
      | DIALOG          | KEYS                             |
      | commits to beam | space down down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                 |
      | old    | git checkout -b parent main                                                                             |
      | parent | git cherry-pick {{ sha-initial 'commit 1' }}                                                            |
      |        | git cherry-pick {{ sha-initial 'commit 4' }}                                                            |
      |        | git checkout old                                                                                        |
      | old    | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1' }}^ {{ sha-initial 'commit 1' }} |
      |        | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 4' }}^ {{ sha-initial 'commit 4' }} |
      |        | git merge --no-edit --ff parent                                                                         |
      |        | git push --force-with-lease --force-if-includes                                                         |
      |        | git checkout parent                                                                                     |
    And this lineage exists now
      """
      main
        parent
          old
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | parent | local         | commit 1                       |
      |        |               | commit 4                       |
      | old    | local, origin | commit 1                       |
      |        |               | commit 2                       |
      |        |               | commit 3                       |
      |        |               | Merge branch 'parent' into old |
      |        | origin        | commit 1                       |
      |        |               | commit 4                       |

  Scenario: sync
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | parent | git fetch --prune --tags  |
      |        | git reset --soft main --  |
      |        | git commit -m "commit 1"  |
      |        | git push -u origin parent |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | parent | local, origin | commit 1                       |
      | old    | local, origin | commit 1                       |
      |        |               | commit 2                       |
      |        |               | commit 3                       |
      |        |               | commit 1                       |
      |        |               | commit 4                       |
      |        |               | Merge branch 'parent' into old |

  Scenario: sync after amending the beamed commit
    And I amend this commit
      | BRANCH | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT    |
      | parent | local    | commit 4b | file_4    | amended content |
    And the current branch is "old"
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                         |
      | old    | git fetch --prune --tags        |
      |        | git checkout parent             |
      | parent | git reset --soft main --        |
      |        | git commit -m "commit 1"        |
      |        | git push -u origin parent       |
      |        | git checkout old                |
      | old    | git merge --no-edit --ff parent |
      |        | git reset --soft parent --      |
      |        | git commit -m "commit 1"        |
      |        | git push --force-with-lease     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | parent | local, origin | commit 1 |
      | old    | local, origin | commit 1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | parent | git checkout old                                |
      | old    | git reset --hard {{ sha 'commit 4' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch -D parent                            |
    And the initial lineage exists now
    And the initial commits exist now
