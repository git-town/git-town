Feature: append a new feature branch in a clean workspace using the "compress" sync strategy with additional commits on the parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE            |
      | feature | local, origin | already compressed |
      | main    | local, origin | new commit         |
    And the current branch is "feature"
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "already compressed"      |
      |         | git push --force-with-lease             |
      |         | git checkout -b new                     |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE            |
      | main    | local, origin | new commit         |
      | feature | local, origin | already compressed |
    And this lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                    |
      | new     | git checkout feature                                       |
      | feature | git reset --hard {{ sha-before-run 'already compressed' }} |
      |         | git push --force-with-lease --force-if-includes            |
      |         | git branch -D new                                          |
    And the initial commits exist now
    And the initial lineage exists now
