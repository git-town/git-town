Feature: sync the entire stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | PARENT | LOCATIONS     |
      | alpha      | feature   | main   | local, origin |
      | beta       | feature   | alpha  | local, origin |
      | gamma      | feature   | beta   | local, origin |
      | one        | feature   | main   | local, origin |
      | two        | feature   | one    | local, origin |
      | production | perennial |        | local, origin |
      | qa         | perennial |        | local, origin |
      | observed   | observed  |        | local, origin |
      | parked     | parked    | main   | local, origin |
    And the commits
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | origin        | main commit              |
      | alpha      | local, origin | alpha commit             |
      | beta       | local, origin | beta commit              |
      | gamma      | local, origin | gamma commit             |
      | one        | local, origin | one commit               |
      | two        | local, origin | two commit               |
      | observed   | local         | local observed commit    |
      |            | origin        | origin observed commit   |
      | parked     | local         | local parked commit      |
      |            | origin        | origin parked commit     |
      | production | local         | local production commit  |
      |            | origin        | origin production commit |
      | qa         | local         | qa local commit          |
      |            | origin        | qa origin commit         |
    And the current branch is "alpha"
    And Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town sync --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main                          |
      |        | git checkout alpha                              |
      | alpha  | git rebase main                                 |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git rebase alpha                                |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma                              |
      | gamma  | git rebase beta                                 |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout alpha                              |
    And the current branch is still "alpha"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | local, origin | main commit              |
      | alpha      | local, origin | main commit              |
      |            |               | alpha commit             |
      | beta       | local, origin | main commit              |
      |            |               | alpha commit             |
      |            |               | beta commit              |
      | gamma      | local, origin | main commit              |
      |            |               | alpha commit             |
      |            |               | beta commit              |
      |            |               | gamma commit             |
      | observed   | local         | local observed commit    |
      |            | origin        | origin observed commit   |
      | one        | local, origin | one commit               |
      | parked     | local         | local parked commit      |
      |            | origin        | origin parked commit     |
      | production | local         | local production commit  |
      |            | origin        | origin production commit |
      | qa         | local         | qa local commit          |
      |            | origin        | qa origin commit         |
      | two        | local, origin | two commit               |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git reset --hard {{ sha 'alpha commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'beta commit' }}        |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma                              |
      | gamma  | git reset --hard {{ sha 'gamma commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git checkout alpha                              |
    And the current branch is still "alpha"
    And the initial commits exist
    And the initial branches and lineage exist
