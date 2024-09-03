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
    When I run "git-town sync --all --detached"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                               |
      | alpha    | git fetch --prune --tags              |
      |          | git merge --no-edit --ff origin/alpha |
      |          | git merge --no-edit --ff main         |
      |          | git checkout beta                     |
      | beta     | git merge --no-edit --ff origin/beta  |
      |          | git merge --no-edit --ff main         |
      |          | git checkout observed                 |
      | observed | git rebase origin/observed            |
      |          | git checkout alpha                    |
      | alpha    | git push --tags                       |
    And the current branch is still "alpha"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | origin        | main commit              |
      | alpha      | local, origin | alpha commit             |
      | beta       | local, origin | beta commit              |
      | observed   | local, origin | origin observed commit   |
      |            | local         | local observed commit    |
      | parked     | local         | local parked commit      |
      |            | origin        | origin parked commit     |
      | production | local         | local production commit  |
      |            | origin        | origin production commit |
      | qa         | local         | qa local commit          |
      |            | origin        | qa origin commit         |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                            |
      | alpha    | git checkout observed                              |
      | observed | git reset --hard {{ sha 'local observed commit' }} |
      |          | git checkout alpha                                 |
    And the current branch is still "alpha"
    And these commits exist now
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
    And the initial branches and lineage exist
