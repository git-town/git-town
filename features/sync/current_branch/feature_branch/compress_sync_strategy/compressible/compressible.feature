Feature: sync a feature branch with multiple commits using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE        | FILE NAME  | FILE CONTENT |
      | alpha  | local, origin | alpha commit 1 | alpha_file | content 1    |
      |        |               | alpha commit 2 | alpha_file | content 2    |
      | beta   | local, origin | beta commit 1  | beta_file  | content 3    |
      |        |               | beta commit 2  | beta_file  | content 4    |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "beta"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                        |
      | beta   | git fetch --prune --tags       |
      |        | git checkout alpha             |
      | alpha  | git reset --soft main --       |
      |        | git commit -m "alpha commit 1" |
      |        | git push --force-with-lease    |
      |        | git checkout beta              |
      | beta   | git merge --no-edit --ff alpha |
      |        | git reset --soft alpha --      |
      |        | git commit -m "beta commit 1"  |
      |        | git push --force-with-lease    |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        | FILE NAME  | FILE CONTENT |
      | alpha  | local, origin | alpha commit 1 | alpha_file | content 2    |
      | beta   | local, origin | beta commit 1  | beta_file  | content 4    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                             |
      | beta   | git checkout alpha                                  |
      | alpha  | git reset --hard {{ sha-initial 'alpha commit 2' }} |
      |        | git push --force-with-lease --force-if-includes     |
      |        | git checkout beta                                   |
      | beta   | git reset --hard {{ sha-initial 'beta commit 2' }}  |
      |        | git push --force-with-lease --force-if-includes     |
    And the initial branches and lineage exist now
    And the initial commits exist now
