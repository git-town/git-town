Feature: sync a feature branch with new commits on the main branch in detached mode

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
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "beta"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                        |
      | beta   | git fetch --prune --tags       |
      |        | git checkout alpha             |
      | alpha  | git checkout beta              |
      | beta   | git merge --no-edit --ff alpha |
      |        | git reset --soft alpha --      |
      |        | git commit -m "beta commit"    |
      |        | git push --force-with-lease    |
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                          |
      | beta   | git reset --hard {{ sha-initial 'beta commit' }} |
      |        | git push --force-with-lease --force-if-includes  |
    And the initial branches and lineage exist now
    And the initial commits exist now
