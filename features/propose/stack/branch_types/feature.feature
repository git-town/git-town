Feature: propose an entire stack

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
      | branch-3 | local, origin | commit 3 |
    And the current branch is "branch-2"
    And tool "open" is installed
    When I run "git-town propose --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                        |
      | branch-2 | git fetch --prune --tags                                                       |
      |          | git checkout branch-1                                                          |
      | branch-1 | git checkout branch-2                                                          |
      | branch-2 | git merge --no-edit --ff branch-1                                              |
      |          | git push                                                                       |
      |          | git checkout branch-3                                                          |
      | branch-3 | git merge --no-edit --ff branch-2                                              |
      |          | git push                                                                       |
      |          | Finding proposal from branch-1 into main ...                                   |
      |          | open https://github.com/git-town/git-town/compare/branch-1?expand=1            |
      |          | Finding proposal from branch-2 into branch-1 ...                               |
      |          | open https://github.com/git-town/git-town/compare/branch-1...branch-2?expand=1 |
      |          | Finding proposal from branch-3 into branch-2 ...                               |
      |          | open https://github.com/git-town/git-town/compare/branch-2...branch-3?expand=1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | git checkout branch-2                           |
      | branch-2 | git reset --hard {{ sha 'commit 2' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the initial lineage exists now
    And the initial branches exist now
