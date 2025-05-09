Feature: sync a feature branch that is already compressed using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | alpha  | local, origin | alpha commit 1 |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | beta   | local, origin | beta commit 1 |
    And wait 1 second to ensure new Git timestamps
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "beta"
    When I run "git-town sync --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                        |
      | beta   | git fetch --prune --tags       |
      |        | git checkout alpha             |
      | alpha  | git checkout beta              |
      | beta   | git merge --no-edit --ff alpha |
      |        | git reset --soft alpha         |
      |        | git commit -m "beta commit 1"  |
      |        | git push --force-with-lease    |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | alpha  | local, origin | alpha commit 1 |
      | beta   | local, origin | beta commit 1  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                               |
      | beta   | git reset --hard {{ sha-before-run 'beta commit 1' }} |
      |        | git push --force-with-lease --force-if-includes       |
    And the initial commits exist now
    And the initial branches and lineage exist now
