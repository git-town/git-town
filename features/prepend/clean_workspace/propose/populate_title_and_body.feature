Feature: propose a newly prepended branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE          |
      | old    | local, origin | old commit       |
      |        |               | unrelated commit |
    And the current branch is "old"
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    And a proposal for this branch does not exist
    When I run "git-town prepend parent --beam --propose --title='proposal title' --body='proposal body'" and enter into the dialog:
      | DIALOG                    | KEYS             |
      | select "unrelated commit" | down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                   |
      | old    | git fetch --prune --tags                                                                                  |
      | <none> | Looking for proposal online ... ok                                                                        |
      | old    | git checkout main                                                                                         |
      | main   | git rebase origin/main --no-update-refs                                                                   |
      |        | git checkout old                                                                                          |
      | old    | git rebase main --no-update-refs                                                                          |
      |        | git push --force-with-lease --force-if-includes                                                           |
      |        | git checkout -b parent main                                                                               |
      | parent | git cherry-pick {{ sha-before-run 'unrelated commit' }}                                                   |
      |        | git checkout old                                                                                          |
      | old    | git rebase parent --no-update-refs                                                                        |
      |        | git push --force-with-lease --force-if-includes                                                           |
      |        | git checkout parent                                                                                       |
      | parent | git push -u origin parent                                                                                 |
      | <none> | open https://github.com/git-town/git-town/compare/parent?expand=1&title=proposal+title&body=proposal+body |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent?expand=1
      """
    And the current branch is now "parent"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE          |
      | old    | local, origin | old commit       |
      | parent | local, origin | unrelated commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the current branch is now "old"
    And the initial commits exist now
    And the initial lineage exists now
