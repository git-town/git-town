Feature: propose an entire stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the current branch is "branch-2"
    And a proposal for this branch exists at "https://github.com/git-town/git-town/pull/2"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose --stack"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                        |
      | branch-2 | git fetch --prune --tags                                                       |
      |          | git checkout branch-1                                                          |
      | branch-1 | git merge --no-edit --ff main                                                  |
      |          | git merge --no-edit --ff origin/branch-1                                       |
      |          | git push                                                                       |
      |          | git checkout branch-2                                                          |
      | branch-2 | git merge --no-edit --ff branch-2                                              |
      |          | git merge --no-edit --ff origin/branch-2                                       |
      |          | git push                                                                       |
      |          | git checkout branch-3                                                          |
      | branch-3 | git merge --no-edit --ff branch-3                                              |
      |          | git merge --no-edit --ff origin/branch-3                                       |
      |          | git push                                                                       |
      | (none)   | Looking for proposal online ... ok                                             |
      |          | open https://github.com/git-town/git-town/compare/branch-1?expand=1            |
      |          | open https://github.com/git-town/git-town/pull/2                               |
      |          | open https://github.com/git-town/git-town/compare/branch-2...branch-3?expand=1 |

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
