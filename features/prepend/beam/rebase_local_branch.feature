@messyoutput
Feature: prepend a branch to a local feature branch using the "rebase" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS |
      | old  | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | old    | local    | commit 1 |
      | old    | local    | commit 2 |
    And the current branch is "old"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town prepend parent --beam" and enter into the dialog:
      | DIALOG          | KEYS        |
      | select commit 1 | space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | old    | git checkout -b parent main                     |
      | parent | git cherry-pick {{ sha-before-run 'commit 1' }} |
      |        | git checkout old                                |
      | old    | git -c rebase.updateRefs=false rebase parent    |
      |        | git checkout parent                             |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE  |
      | old    | local    | commit 2 |
      | parent | local    | commit 1 |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | parent | git checkout old                      |
      | old    | git reset --hard {{ sha 'commit 2' }} |
      |        | git branch -D parent                  |
    And the initial commits exist now
    And the initial lineage exists now
