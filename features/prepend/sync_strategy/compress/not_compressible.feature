Feature: prepend a branch to a feature branch that is already compressed in a clean workspace using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-1 | local, origin | branch-1 commit |
      | branch-2 | local, origin | branch-2 commit |
    And the current branch is "branch-2"
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town prepend branch-1a"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                            |
      | branch-2 | git fetch --prune --tags           |
      |          | git checkout branch-1              |
      | branch-1 | git checkout branch-2              |
      | branch-2 | git merge --no-edit --ff branch-1  |
      |          | git reset --soft branch-1          |
      |          | git commit -m "branch-2 commit"    |
      |          | git push --force-with-lease        |
      |          | git checkout -b branch-1a branch-1 |
    And the initial commits exist now
    And this lineage exists now
      | BRANCH    | PARENT    |
      | branch-1  | main      |
      | branch-1a | branch-1  |
      | branch-2  | branch-1a |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                              |
      | branch-1a | git checkout branch-2                                |
      | branch-2  | git reset --hard {{ sha-initial 'branch-2 commit' }} |
      |           | git push --force-with-lease --force-if-includes      |
      |           | git branch -D branch-1a                              |
    And the initial commits exist now
    And the initial lineage exists now
