Feature: sync the current feature branch with a tracking branch in detached mode with updates on the tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | alpha  | origin   | alpha commit |
      | beta   | origin   | beta commit  |
    And the current branch is "beta"
    When I run "git-town sync --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git push                              |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                              |
      | alpha  | local, origin | alpha commit                                         |
      | beta   | local, origin | beta commit                                          |
      |        |               | Merge remote-tracking branch 'origin/beta' into beta |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                   |
      | beta   | git reset --hard {{ sha 'initial commit' }}                               |
      |        | git push --force-with-lease origin {{ sha-in-origin 'beta commit' }}:beta |
      |        | git checkout alpha                                                        |
      | alpha  | git reset --hard {{ sha 'initial commit' }}                               |
      |        | git checkout beta                                                         |
    And the initial branches and lineage exist now
    And the initial commits exist now
