Feature: proposing a branch whose parent was shipped and the local branch deleted manually

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin ships the "parent" branch using the "squash-merge" ship-strategy
    And I ran "git branch -d parent"
    And the current branch is "child"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                          |
      | child  | git fetch --prune --tags                                         |
      | (none) | open https://github.com/git-town/git-town/compare/child?expand=1 |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And the initial lineage exists now
