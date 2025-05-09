Feature: sync a feature branch with new commits on the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | new commit   |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
    And wait 1 second to ensure new Git timestamps
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "beta"
    When I run "git-town sync --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff main         |
      |        | git merge --no-edit --ff origin/alpha |
      |        | git reset --soft main                 |
      |        | git commit -m "alpha commit"          |
      |        | git push --force-with-lease           |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git reset --soft alpha                |
      |        | git commit -m "beta commit"           |
      |        | git push --force-with-lease           |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | new commit   |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git checkout alpha                                   |
      | alpha  | git reset --hard {{ sha-before-run 'alpha commit' }} |
      |        | git push --force-with-lease --force-if-includes      |
      |        | git checkout beta                                    |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git push --force-with-lease --force-if-includes      |
    And the initial commits exist now
    And the initial branches and lineage exist now
