@messyoutput
Feature: propose a newly prepended branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE          |
      | existing | local, origin | existing commit  |
      |          |               | unrelated commit |
    And the current branch is "existing"
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    And a proposal for this branch does not exist
    When I run "git-town prepend new --beam --propose" and enter into the dialog:
      | DIALOG                    | KEYS             |
      | select "unrelated commit" | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | existing | git fetch --prune --tags                                       |
      | <none>   | Looking for proposal online ... ok                             |
      | existing | git checkout main                                              |
      | main     | git rebase origin/main --no-update-refs                        |
      |          | git checkout existing                                          |
      | existing | git rebase main --no-update-refs                               |
      |          | git push --force-with-lease --force-if-includes                |
      |          | git checkout -b new main                                       |
      | new      | git cherry-pick {{ sha-before-run 'unrelated commit' }}        |
      |          | git checkout existing                                          |
      | existing | git rebase new --no-update-refs                                |
      |          | git push --force-with-lease --force-if-includes                |
      |          | git checkout new                                               |
      | new      | git push -u origin new                                         |
      | <none>   | open https://github.com/git-town/git-town/compare/new?expand=1 |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/new?expand=1
      """
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE          |
      | existing | local, origin | existing commit  |
      | new      | local, origin | unrelated commit |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | new    |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | new      | git checkout existing                           |
      | existing | git reset --hard {{ sha 'unrelated commit' }}   |
      |          | git push --force-with-lease --force-if-includes |
      |          | git branch -D new                               |
      |          | git push origin :new                            |
    And the current branch is now "existing"
    And the initial commits exist now
    And the initial lineage exists now
