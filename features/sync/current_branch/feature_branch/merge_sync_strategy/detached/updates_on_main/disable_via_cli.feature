Feature: sync the current feature branch with a tracking branch in detached mode with updates on main

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And Git setting "git-town.detached" is "true"
    And the current branch is "beta"
    When I run "git-town sync --no-detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                        |
      | beta   | git fetch --prune --tags       |
      |        | git checkout alpha             |
      | alpha  | git merge --no-edit --ff main  |
      |        | git push                       |
      |        | git checkout beta              |
      | beta   | git merge --no-edit --ff alpha |
      |        | git push                       |
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'initial commit` }}     |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'initial commit` }}     |
      |        | git push --force-with-lease --force-if-includes |
    And the initial branches and lineage exist now
    And the initial commits exist now
