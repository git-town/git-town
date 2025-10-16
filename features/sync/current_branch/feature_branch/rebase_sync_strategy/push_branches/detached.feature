Feature: detached syncing a stacked feature branch using --no-push

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
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "beta"
    When I run "git-town sync --no-push --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                       |
      | beta   | git fetch --prune --tags                                                      |
      |        | git checkout alpha                                                            |
      | alpha  | git -c rebase.updateRefs=false rebase origin/alpha                            |
      |        | git checkout beta                                                             |
      | beta   | git -c rebase.updateRefs=false rebase origin/beta                             |
      |        | git -c rebase.updateRefs=false rebase --onto alpha {{ sha 'initial commit' }} |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | local         | local main commit   |
      |        | origin        | origin main commit  |
      | alpha  | local, origin | origin alpha commit |
      |        | local         | local alpha commit  |
      | beta   | local         | origin beta commit  |
      |        |               | local beta commit   |
      |        | origin        | origin beta commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'local alpha commit' }} |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'local beta commit' }}  |
    And the initial branches and lineage exist now
    And the initial commits exist now
