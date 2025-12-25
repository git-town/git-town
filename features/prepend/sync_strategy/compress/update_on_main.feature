Feature: prepend a branch to a feature branch in a clean workspace using the "compress" sync strategy when there are new commits on main

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
      | main     | local, origin | main commit     |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "branch-2"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town prepend branch-1a"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                            |
      | branch-2 | git fetch --prune --tags           |
      |          | git checkout branch-1              |
      | branch-1 | git merge --no-edit --ff main      |
      |          | git reset --soft main --           |
      |          | git commit -m "branch-1 commit"    |
      |          | git push --force-with-lease        |
      |          | git checkout branch-2              |
      | branch-2 | git merge --no-edit --ff branch-1  |
      |          | git reset --soft branch-1 --       |
      |          | git commit -m "branch-2 commit"    |
      |          | git push --force-with-lease        |
      |          | git checkout -b branch-1a branch-1 |
    And this lineage exists now
      """
      main
        branch-1
          branch-1a
            branch-2
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                              |
      | branch-1a | git checkout branch-1                                |
      | branch-1  | git reset --hard {{ sha-initial 'branch-1 commit' }} |
      |           | git push --force-with-lease --force-if-includes      |
      |           | git checkout branch-2                                |
      | branch-2  | git reset --hard {{ sha-initial 'branch-2 commit' }} |
      |           | git push --force-with-lease --force-if-includes      |
      |           | git branch -D branch-1a                              |
    And the initial lineage exists now
    And the initial commits exist now
