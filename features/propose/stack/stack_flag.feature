Feature: propose an entire stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
      | branch-4 | feature | branch-3 | local, origin |
    And the current branch is "branch-3"
    And a proposal for this branch exists at "https://github.com/git-town/git-town/pull/3"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose --stack"

  @this
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
      | (none) | open https://github.com/git-town/git-town/compare/child?expand=1 |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
    And Git Town runs the commands
      | BRANCH | COMMAND                                                          |
      | child  | git fetch --prune --tags                                         |
      |        | git checkout main                                                |
      | main   | git rebase origin/main --no-update-refs                          |
      |        | git checkout child                                               |
      | child  | git merge --no-edit --ff main                                    |
      |        | git merge --no-edit --ff origin/child                            |
      |        | git push                                                         |
      | (none) | open https://github.com/git-town/git-town/compare/child?expand=1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git reset --hard {{ sha 'child commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git checkout child                              |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And the initial lineage exists now
