Feature: append a new feature branch in a clean workspace using the "compress" sync strategy with new commits on the parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE            |
      | feature | local, origin | already compressed |
      | main    | local, origin | new commit         |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "feature"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git merge --no-edit --ff main      |
      |         | git reset --soft main --           |
      |         | git commit -m "already compressed" |
      |         | git push --force-with-lease        |
      |         | git checkout -b new                |
    And this lineage exists now
      """
      main
        feature
          new
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                 |
      | new     | git checkout feature                                    |
      | feature | git reset --hard {{ sha-initial 'already compressed' }} |
      |         | git push --force-with-lease --force-if-includes         |
      |         | git branch -D new                                       |
    And the initial lineage exists now
    And the initial commits exist now
