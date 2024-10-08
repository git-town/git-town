Feature: sync all feature branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | PARENT | LOCATIONS     |
      | alpha      | feature   | main   | local, origin |
      | beta       | feature   | main   | local, origin |
      | production | perennial |        | local, origin |
      | qa         | perennial |        | local, origin |
      | observed   | observed  |        | local, origin |
      | parked     | parked    | main   | local, origin |
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

  Scenario: with "merge" sync-feature strategy
    When I run "git-town sync --all"
    Then it runs the commands
      | BRANCH     | COMMAND                                       |
      | alpha      | git fetch --prune --tags                      |
      |            | git checkout main                             |
      | main       | git rebase origin/main --no-update-refs       |
      |            | git checkout alpha                            |
      | alpha      | git merge --no-edit --ff origin/alpha         |
      |            | git merge --no-edit --ff main                 |
      |            | git push                                      |
      |            | git checkout beta                             |
      | beta       | git merge --no-edit --ff origin/beta          |
      |            | git merge --no-edit --ff main                 |
      |            | git push                                      |
      |            | git checkout observed                         |
      | observed   | git rebase origin/observed --no-update-refs   |
      |            | git checkout production                       |
      | production | git rebase origin/production --no-update-refs |
      |            | git push                                      |
      |            | git checkout qa                               |
      | qa         | git rebase origin/qa --no-update-refs         |
      |            | git push                                      |
      |            | git checkout alpha                            |
      | alpha      | git push --tags                               |
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
      | BRANCH     | COMMAND                                         |
      | alpha      | git fetch --prune --tags                        |
      |            | git checkout main                               |
      | main       | git rebase origin/main --no-update-refs         |
      |            | git checkout alpha                              |
      | alpha      | git rebase main --no-update-refs                |
      |            | git push --force-with-lease --force-if-includes |
      |            | git checkout beta                               |
      | beta       | git rebase main --no-update-refs                |
      |            | git push --force-with-lease --force-if-includes |
      |            | git checkout observed                           |
      | observed   | git rebase origin/observed --no-update-refs     |
      |            | git checkout production                         |
      | production | git rebase origin/production --no-update-refs   |
      |            | git push                                        |
      |            | git checkout qa                                 |
      | qa         | git rebase origin/qa --no-update-refs           |
      |            | git push                                        |
      |            | git checkout alpha                              |
      | alpha      | git push --tags                                 |
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
