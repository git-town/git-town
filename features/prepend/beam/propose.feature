@messyoutput
Feature: propose a newly prepended branch

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | parent   | feature | main   | local, origin |
      | existing | feature | parent | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE          |
      | existing | local, origin | existing commit  |
      |          |               | unrelated commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "existing"
    And a proposal for this branch does not exist
    And tool "open" is installed
    When I run "git-town prepend new --beam --propose" and enter into the dialog:
      | DIALOG          | KEYS             |
      | commits to beam | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                                 |
      |          | Finding proposal from existing into parent ... ok                                                                       |
      | existing | git checkout -b new parent                                                                                              |
      | new      | git cherry-pick {{ sha-initial 'unrelated commit' }}                                                                    |
      |          | git checkout existing                                                                                                   |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'unrelated commit' }}^ {{ sha-initial 'unrelated commit' }} |
      |          | git -c rebase.updateRefs=false rebase new                                                                               |
      |          | git push --force-with-lease --force-if-includes                                                                         |
      |          | git checkout new                                                                                                        |
      | new      | git push -u origin new                                                                                                  |
      |          | Finding proposal from new into parent ... ok                                                                            |
      |          | open https://github.com/git-town/git-town/compare/parent...new?expand=1                                                 |
    And this lineage exists now
      """
      main
        parent
          new
            existing
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE          |
      | new      | local, origin | unrelated commit |
      | existing | local, origin | existing commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | new      | git checkout existing                           |
      | existing | git reset --hard {{ sha 'unrelated commit' }}   |
      |          | git push --force-with-lease --force-if-includes |
      |          | git branch -D new                               |
      |          | git push origin :new                            |
    And the initial lineage exists now
    And the initial commits exist now
