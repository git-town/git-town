Feature: auto-propose new branches

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And Git setting "git-town.share-new-branches" is "propose"
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                        |
      | main   | git fetch --prune --tags                                       |
      |        | git -c rebase.updateRefs=false rebase origin/main              |
      |        | git checkout -b new                                            |
      | new    | git push -u origin new                                         |
      | (none) | open https://github.com/git-town/git-town/compare/new?expand=1 |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | new    | git checkout main                           |
      | main   | git reset --hard {{ sha 'initial commit' }} |
      |        | git branch -D new                           |
      |        | git push origin :new                        |
    And the initial commits exist now
    And the initial branches and lineage exist now
