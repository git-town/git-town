@smoke
Feature: sync the current feature branch with a tracking branch using the "merge" sync-feature strategy

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
    When I run "git-town sync --detached"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git merge --no-edit --ff main         |
      |        | git push                              |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff origin/beta  |
      |        | git merge --no-edit --ff alpha        |
      |        | git push                              |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                |
      | main   | local         | local main commit                                      |
      |        | origin        | origin main commit                                     |
      | alpha  | local, origin | local alpha commit                                     |
      |        |               | origin alpha commit                                    |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha |
      |        |               | local main commit                                      |
      |        |               | Merge branch 'main' into alpha                         |
      | beta   | local, origin | local beta commit                                      |
      |        |               | origin beta commit                                     |
      |        |               | Merge remote-tracking branch 'origin/beta' into beta   |
      |        |               | local alpha commit                                     |
      |        |               | origin alpha commit                                    |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha |
      |        |               | local main commit                                      |
      |        |               | Merge branch 'main' into alpha                         |
      |        |               | Merge branch 'alpha' into beta                         |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                            |
      | beta   | git checkout alpha                                                                 |
      | alpha  | git reset --hard {{ sha 'local alpha commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin 'origin alpha commit' }}:alpha |
      |        | git checkout beta                                                                  |
      | beta   | git reset --hard {{ sha 'local beta commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin 'origin beta commit' }}:beta   |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial branches and lineage exist now
