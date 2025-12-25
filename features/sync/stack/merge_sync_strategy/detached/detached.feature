Feature: detached sync of the entire stack

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
    And the current branch is "alpha"
    When I run "git-town sync --stack --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | alpha  | git fetch --prune --tags              |
      |        | git merge --no-edit --ff origin/alpha |
      |        | git push                              |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git push                              |
      |        | git checkout alpha                    |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                |
      | main   | local         | local main commit                                      |
      |        | origin        | origin main commit                                     |
      | alpha  | local, origin | local alpha commit                                     |
      |        |               | origin alpha commit                                    |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha |
      | beta   | local, origin | local beta commit                                      |
      |        |               | Merge branch 'alpha' into beta                         |
      |        |               | origin beta commit                                     |
      |        |               | Merge remote-tracking branch 'origin/beta' into beta   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                    |
      | alpha  | git reset --hard {{ sha-initial 'local alpha commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin alpha commit' }}:alpha |
      |        | git checkout beta                                                                          |
      | beta   | git reset --hard {{ sha-initial 'local beta commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin beta commit' }}:beta   |
      |        | git checkout alpha                                                                         |
