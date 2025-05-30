Feature: detached sync a grandchild feature branch using the "compress" strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE             |
      | main   | local    | local main commit   |
      |        | origin   | origin main commit  |
      | alpha  | local    | local alpha commit  |
      |        | origin   | origin alpha commit |
      | beta   | local    | local beta commit   |
      |        | origin   | origin beta commit  |
    And wait 1 second to ensure new Git timestamps
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "beta"
    When I run "git-town sync --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git reset --soft main                 |
      |        | git commit -m "local alpha commit"    |
      |        | git push --force-with-lease           |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git reset --soft alpha                |
      |        | git commit -m "local beta commit"     |
      |        | git push --force-with-lease           |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local         | local main commit  |
      |        | origin        | origin main commit |
      | alpha  | local, origin | local alpha commit |
      |        | origin        | local main commit  |
      | beta   | local, origin | local beta commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                    |
      | beta   | git checkout alpha                                                                         |
      | alpha  | git reset --hard {{ sha-initial 'local alpha commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin alpha commit' }}:alpha |
      |        | git checkout beta                                                                          |
      | beta   | git reset --hard {{ sha-initial 'local beta commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin beta commit' }}:beta   |
    And the initial commits exist now
    And the initial branches and lineage exist now
