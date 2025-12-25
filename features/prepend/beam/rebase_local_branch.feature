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
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "old"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town prepend parent --beam" and enter into the dialog:
      | DIALOG          | KEYS        |
      | commits to beam | space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                 |
      | old    | git checkout -b parent main                                                                             |
      | parent | git cherry-pick {{ sha-initial 'commit 1' }}                                                            |
      |        | git checkout old                                                                                        |
      | old    | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 1' }}^ {{ sha-initial 'commit 1' }} |
      |        | git -c rebase.updateRefs=false rebase parent                                                            |
      |        | git checkout parent                                                                                     |
    And this lineage exists now
      """
      main
        parent
          old
      """
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE  |
      | parent | local    | commit 1 |
      | old    | local    | commit 2 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | parent | git checkout old                      |
      | old    | git reset --hard {{ sha 'commit 2' }} |
      |        | git branch -D parent                  |
    And the initial lineage exists now
    And the initial commits exist now
