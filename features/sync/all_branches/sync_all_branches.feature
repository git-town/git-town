Feature: sync all feature branches

  Background:
    Given the feature branches "alpha" and "beta"
    And the perennial branches "production" and "qa"
    And an observed branch "observed"
    And a parked branch "parked"
    And the commits
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | origin        | main commit              |
      | alpha      | local, origin | alpha commit             |
      | beta       | local, origin | beta commit              |
      | observed   | local         | local observed commit    |
      |            | origin        | origin observed commit   |
      | parked     | local         | local parked commit      |
      |            | origin        | origin parked commit     |
      | production | local         | local production commit  |
      |            | origin        | origin production commit |
      | qa         | local         | qa local commit          |
      |            | origin        | qa origin commit         |
    And the current branch is "alpha"

  @this
  Scenario: with "merge" sync-feature strategy
    When I run "git-town sync --all"
    Then it runs the commands
      | BRANCH     | COMMAND                          |
      | alpha      | git fetch --prune --tags         |
      |            | git checkout main                |
      | main       | git rebase origin/main           |
      |            | git checkout alpha               |
      | alpha      | git merge --no-edit origin/alpha |
      |            | git merge --no-edit main         |
      |            | git push                         |
      |            | git checkout beta                |
      | beta       | git merge --no-edit origin/beta  |
      |            | git merge --no-edit main         |
      |            | git push                         |
      |            | git checkout observed            |
      | observed   | git rebase origin/observed       |
      |            | git checkout production          |
      | production | git rebase origin/production     |
      |            | git push                         |
      |            | git checkout qa                  |
      | qa         | git rebase origin/qa             |
      |            | git push                         |
      |            | git checkout alpha               |
      | alpha      | git push --tags                  |
    And the current branch is still "alpha"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                        |
      | main       | local, origin | main commit                    |
      | alpha      | local, origin | alpha commit                   |
      |            |               | main commit                    |
      |            |               | Merge branch 'main' into alpha |
      | beta       | local, origin | beta commit                    |
      |            |               | main commit                    |
      |            |               | Merge branch 'main' into beta  |
      | observed   | local, origin | origin observed commit         |
      |            | local         | local observed commit          |
      | parked     | local         | local parked commit            |
      |            | origin        | origin parked commit           |
      | production | local, origin | origin production commit       |
      |            |               | local production commit        |
      | qa         | local, origin | qa origin commit               |
      |            |               | qa local commit                |

  Scenario: with "rebase" sync-feature strategy
    Given Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town sync --all"
    Then it runs the commands
      | BRANCH     | COMMAND                      |
      | alpha      | git fetch --prune --tags     |
      |            | git checkout main            |
      | main       | git rebase origin/main       |
      |            | git checkout alpha           |
      | alpha      | git rebase origin/alpha      |
      |            | git rebase main              |
      |            | git push --force-with-lease  |
      |            | git checkout beta            |
      | beta       | git rebase origin/beta       |
      |            | git rebase main              |
      |            | git push --force-with-lease  |
      |            | git checkout observed        |
      | observed   | git rebase origin/observed   |
      |            | git checkout production      |
      | production | git rebase origin/production |
      |            | git push                     |
      |            | git checkout qa              |
      | qa         | git rebase origin/qa         |
      |            | git push                     |
      |            | git checkout alpha           |
      | alpha      | git push --tags              |
    And the current branch is still "alpha"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | local, origin | main commit              |
      | alpha      | local, origin | main commit              |
      |            |               | alpha commit             |
      | beta       | local, origin | main commit              |
      |            |               | beta commit              |
      | observed   | local, origin | origin observed commit   |
      |            | local         | local observed commit    |
      | parked     | local         | local parked commit      |
      |            | origin        | origin parked commit     |
      | production | local, origin | origin production commit |
      |            |               | local production commit  |
      | qa         | local, origin | qa origin commit         |
      |            |               | qa local commit          |
