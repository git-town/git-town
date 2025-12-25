Feature: auto-propose new branches

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    And Git setting "git-town.share-new-branches" is "propose"
    And the current branch is "old"
    And tool "open" is installed
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                        |
      | old    | git fetch --prune --tags                                       |
      |        | git checkout -b new main                                       |
      | new    | git push -u origin new                                         |
      |        | open https://github.com/git-town/git-town/compare/new?expand=1 |
    And this lineage exists now
      """
      main
        new
          old
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | new    | git checkout old     |
      | old    | git branch -D new    |
      |        | git push origin :new |
    And the initial lineage exists now
    And the initial commits exist now
