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

  Scenario: with "merge" feature sync strategy
    When I run "git-town sync --all"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                 |
      | alpha      | git fetch --prune --tags                                |
      |            | git checkout main                                       |
      | main       | git -c rebase.updateRefs=false rebase origin/main       |
      |            | git checkout alpha                                      |
      | alpha      | git merge --no-edit --ff main                           |
      |            | git push                                                |
      |            | git checkout beta                                       |
      | beta       | git merge --no-edit --ff main                           |
      |            | git push                                                |
      |            | git checkout observed                                   |
      | observed   | git -c rebase.updateRefs=false rebase origin/observed   |
      |            | git checkout production                                 |
      | production | git -c rebase.updateRefs=false rebase origin/production |
      |            | git push                                                |
      |            | git checkout qa                                         |
      | qa         | git -c rebase.updateRefs=false rebase origin/qa         |
      |            | git push                                                |
      |            | git checkout alpha                                      |
      | alpha      | git push --tags                                         |
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                        |
      | main       | local, origin | main commit                    |
      | alpha      | local, origin | alpha commit                   |
      |            |               | Merge branch 'main' into alpha |
      | beta       | local, origin | beta commit                    |
      |            |               | Merge branch 'main' into beta  |
      | parked     | local         | local parked commit            |
      |            | origin        | origin parked commit           |
      | observed   | local, origin | origin observed commit         |
      |            | local         | local observed commit          |
      | production | local, origin | origin production commit       |
      |            |               | local production commit        |
      | qa         | local, origin | qa origin commit               |
      |            |               | qa local commit                |

  Scenario: with "rebase" feature sync strategy
    Given Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync --all"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                                      |
      | alpha      | git fetch --prune --tags                                                     |
      |            | git checkout main                                                            |
      | main       | git -c rebase.updateRefs=false rebase origin/main                            |
      |            | git checkout alpha                                                           |
      | alpha      | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |            | git push --force-with-lease --force-if-includes                              |
      |            | git checkout beta                                                            |
      | beta       | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |            | git push --force-with-lease --force-if-includes                              |
      |            | git checkout observed                                                        |
      | observed   | git -c rebase.updateRefs=false rebase origin/observed                        |
      |            | git checkout production                                                      |
      | production | git -c rebase.updateRefs=false rebase origin/production                      |
      |            | git push                                                                     |
      |            | git checkout qa                                                              |
      | qa         | git -c rebase.updateRefs=false rebase origin/qa                              |
      |            | git push                                                                     |
      |            | git checkout alpha                                                           |
      | alpha      | git push --tags                                                              |
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | local, origin | main commit              |
      | alpha      | local, origin | alpha commit             |
      | beta       | local, origin | beta commit              |
      | parked     | local         | local parked commit      |
      |            | origin        | origin parked commit     |
      | observed   | local, origin | origin observed commit   |
      |            | local         | local observed commit    |
      | production | local, origin | origin production commit |
      |            |               | local production commit  |
      | qa         | local, origin | qa origin commit         |
      |            |               | qa local commit          |
