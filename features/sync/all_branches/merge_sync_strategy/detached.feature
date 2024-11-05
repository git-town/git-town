Feature: sync all feature branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | alpha        | feature      | main   | local, origin |
      | beta         | feature      | main   | local, origin |
      | production   | perennial    |        | local, origin |
      | qa           | perennial    |        | local, origin |
      | observed     | observed     |        | local, origin |
      | contribution | contribution |        | local, origin |
      | parked       | parked       | main   | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE                    |
      | main         | origin        | main commit                |
      | alpha        | local, origin | alpha commit               |
      | beta         | local, origin | beta commit                |
      | contribution | local         | local contribution commit  |
      |              | origin        | origin contribution commit |
      | observed     | local         | local observed commit      |
      |              | origin        | origin observed commit     |
      | parked       | local         | local parked commit        |
      |              | origin        | origin parked commit       |
      | production   | local         | local production commit    |
      |              | origin        | origin production commit   |
      | qa           | local         | qa local commit            |
      |              | origin        | qa origin commit           |
    And the current branch is "alpha"
    When I run "git-town sync --all --detached"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                                         |
      | alpha        | git fetch --prune --tags                        |
      |              | git merge --no-edit --ff main                   |
      |              | git merge --no-edit --ff origin/alpha           |
      |              | git checkout beta                               |
      | beta         | git merge --no-edit --ff main                   |
      |              | git merge --no-edit --ff origin/beta            |
      |              | git checkout contribution                       |
      | contribution | git rebase origin/contribution --no-update-refs |
      |              | git push                                        |
      |              | git checkout observed                           |
      | observed     | git rebase origin/observed --no-update-refs     |
      |              | git checkout alpha                              |
      | alpha        | git push --tags                                 |
    And the current branch is still "alpha"
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE                    |
      | main         | origin        | main commit                |
      | alpha        | local, origin | alpha commit               |
      | beta         | local, origin | beta commit                |
      | contribution | local, origin | origin contribution commit |
      |              |               | local contribution commit  |
      | observed     | local, origin | origin observed commit     |
      |              | local         | local observed commit      |
      | parked       | local         | local parked commit        |
      |              | origin        | origin parked commit       |
      | production   | local         | local production commit    |
      |              | origin        | origin production commit   |
      | qa           | local         | qa local commit            |
      |              | origin        | qa origin commit           |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                                                                          |
      | alpha        | git checkout contribution                                                                        |
      | contribution | git reset --hard {{ sha 'local contribution commit' }}                                           |
      |              | git push --force-with-lease origin {{ sha-in-origin 'origin contribution commit' }}:contribution |
      |              | git checkout observed                                                                            |
      | observed     | git reset --hard {{ sha 'local observed commit' }}                                               |
      |              | git checkout alpha                                                                               |
    And the current branch is still "alpha"
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE                    |
      | main         | origin        | main commit                |
      | alpha        | local, origin | alpha commit               |
      | beta         | local, origin | beta commit                |
      | contribution | local         | local contribution commit  |
      |              | origin        | origin contribution commit |
      | observed     | local         | local observed commit      |
      |              | origin        | origin observed commit     |
      | parked       | local         | local parked commit        |
      |              | origin        | origin parked commit       |
      | production   | local         | local production commit    |
      |              | origin        | origin production commit   |
      | qa           | local         | qa local commit            |
      |              | origin        | qa origin commit           |
    And the initial branches and lineage exist now
