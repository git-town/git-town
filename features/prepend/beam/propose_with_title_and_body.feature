@messyoutput
Feature: propose a newly prepended branch

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE          |
      | existing | local, origin | existing commit  |
      |          |               | unrelated commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "existing"
    And a proposal for this branch does not exist
    And tool "open" is installed
    When I run "git-town prepend new --beam --propose --title='proposal title' --body='proposal body'" and enter into the dialog:
      | DIALOG          | KEYS             |
      | commits to beam | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                                 |
      |          | Looking for proposal online ... ok                                                                                      |
      | existing | git checkout -b new main                                                                                                |
      | new      | git cherry-pick {{ sha-initial 'unrelated commit' }}                                                                    |
      |          | git checkout existing                                                                                                   |
      | existing | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'unrelated commit' }}^ {{ sha-initial 'unrelated commit' }} |
      |          | git -c rebase.updateRefs=false rebase new                                                                               |
      |          | git push --force-with-lease --force-if-includes                                                                         |
      |          | git checkout new                                                                                                        |
      | new      | git push -u origin new                                                                                                  |
      |          | Looking for proposal online ... ok                                                                                      |
      |          | open https://github.com/git-town/git-town/compare/new?expand=1&title=proposal+title&body=proposal+body                  |
    And this lineage exists now
      """
      main
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
