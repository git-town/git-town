Feature: detached syncing a stacked feature branch using --no-push

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the current branch is "beta"
    And the commits
      | BRANCH | LOCATION | MESSAGE             |
      | main   | local    | local main commit   |
      |        | origin   | origin main commit  |
      | alpha  | local    | local alpha commit  |
      |        | origin   | origin alpha commit |
      | beta   | local    | local beta commit   |
      |        | origin   | origin beta commit  |
    And Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town sync --no-push --detached"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                  |
      | beta   | git fetch --prune --tags                 |
      |        | git checkout alpha                       |
      | alpha  | git rebase main --no-update-refs         |
      |        | git rebase origin/alpha --no-update-refs |
      |        | git checkout beta                        |
      | beta   | git rebase alpha --no-update-refs        |
      |        | git rebase origin/beta --no-update-refs  |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | local         | local main commit   |
      |        | origin        | origin main commit  |
      | alpha  | local, origin | origin alpha commit |
      |        | local         | local main commit   |
      |        |               | local alpha commit  |
      | beta   | local, origin | origin beta commit  |
      |        | local         | origin alpha commit |
      |        |               | local main commit   |
      |        |               | local alpha commit  |
      |        |               | local beta commit   |
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'local alpha commit' }} |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'local beta commit' }}  |

    And the current branch is still "beta"
    And the initial commits exist now
    And the initial branches and lineage exist now
