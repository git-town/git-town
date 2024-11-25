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
      |        | git checkout main                                                |
      | main   | git rebase origin/main --no-update-refs                          |
      |        | git checkout child                                               |
      | child  | git merge --no-edit --ff main                                    |
      |        | git merge --no-edit --ff origin/child                            |
      |        | git push                                                         |
      | <none> | open https://github.com/git-town/git-town/compare/child?expand=1 |
    And the current branch is still "child"
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/child?expand=1
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git checkout child                              |
      | child  | git reset --hard {{ sha 'child commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git checkout child                              |
      | child  | git branch -D new                               |
    And the current branch is still "child"
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And the initial lineage exists now
