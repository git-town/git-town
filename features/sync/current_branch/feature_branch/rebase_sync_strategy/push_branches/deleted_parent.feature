Feature: don't push when the parent branch was shipped and pushing is disabled

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
      | beta   | local, origin | beta commit  | beta-file  | beta content  |
    And Git setting "git-town.push-branches" is "false"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "alpha" branch using the "squash-merge" ship-strategy
    And the current branch is "beta"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                 |
      | beta   | git fetch --prune --tags                                |
      |        | git checkout main                                       |
      | main   | git -c rebase.updateRefs=false rebase origin/main       |
      |        | git checkout beta                                       |
      | beta   | git pull                                                |
      |        | git -c rebase.updateRefs=false rebase --onto main alpha |
      |        | git branch -D alpha                                     |
    And the branches are now
      | REPOSITORY    | BRANCHES   |
      | local, origin | main, beta |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | alpha commit |
      | beta   | local         | beta commit  |
      |        | origin        | beta commit  |
