Feature: detached syncing a stacked feature branch using --no-push

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | other | (none)  |        | local, origin |
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
    When I run "git-town sync --no-push --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff main         |
      |        | git merge --no-edit --ff origin/alpha |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                |
      | main   | local         | local main commit                                      |
      |        | origin        | origin main commit                                     |
      | alpha  | local         | local alpha commit                                     |
      |        |               | Merge branch 'main' into alpha                         |
      |        | local, origin | origin alpha commit                                    |
      |        | local         | Merge remote-tracking branch 'origin/alpha' into alpha |
      | beta   | local         | local beta commit                                      |
      |        |               | Merge branch 'alpha' into beta                         |
      |        | local, origin | origin beta commit                                     |
      |        | local         | Merge remote-tracking branch 'origin/beta' into beta   |
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'local alpha commit' }} |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'local beta commit' }}  |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial branches and lineage exist now
