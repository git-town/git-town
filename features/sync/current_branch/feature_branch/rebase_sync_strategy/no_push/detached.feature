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
    And the current branch is "beta"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync --no-push --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                            |
      | beta   | git fetch --prune --tags                           |
      |        | git checkout alpha                                 |
      | alpha  | git -c rebase.updateRefs=false rebase origin/alpha |
      |        | git push --force-with-lease --force-if-includes    |
      |        | git checkout beta                                  |
      | beta   | git -c rebase.updateRefs=false rebase alpha        |
      |        | git -c rebase.updateRefs=false rebase origin/beta  |
      |        | git -c rebase.updateRefs=false rebase alpha        |
      |        | git push --force-with-lease --force-if-includes    |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | local         | local main commit   |
      |        | origin        | origin main commit  |
      | alpha  | local, origin | origin alpha commit |
      |        |               | local alpha commit  |
      | beta   | local, origin | origin beta commit  |
      |        |               | local beta commit   |
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                       |
      | beta   | git checkout alpha                                                                            |
      | alpha  | git reset --hard {{ sha 'local alpha commit' }}                                               |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin alpha commit' }}:alpha |
      |        | git checkout beta                                                                             |
      | beta   | git reset --hard {{ sha 'local beta commit' }}                                                |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin beta commit' }}:beta   |
    And the initial commits exist now
    And the initial branches and lineage exist now
