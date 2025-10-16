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
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                   |
      | alpha        | git fetch --prune --tags                                  |
      |              | git checkout beta                                         |
      | beta         | git checkout contribution                                 |
      | contribution | git -c rebase.updateRefs=false rebase origin/contribution |
      |              | git push                                                  |
      |              | git checkout observed                                     |
      | observed     | git -c rebase.updateRefs=false rebase origin/observed     |
      |              | git checkout alpha                                        |
      | alpha        | git push --tags                                           |
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE                    |
      | main         | origin        | main commit                |
      | alpha        | local, origin | alpha commit               |
      | beta         | local, origin | beta commit                |
      | parked       | local         | local parked commit        |
      |              | origin        | origin parked commit       |
      | contribution | local, origin | origin contribution commit |
      |              |               | local contribution commit  |
      | observed     | local, origin | origin observed commit     |
      |              | local         | local observed commit      |
      | production   | local         | local production commit    |
      |              | origin        | origin production commit   |
      | qa           | local         | qa local commit            |
      |              | origin        | qa origin commit           |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                                                          |
      | alpha        | git checkout contribution                                                                        |
      | contribution | git reset --hard {{ sha 'local contribution commit' }}                                           |
      |              | git push --force-with-lease origin {{ sha-in-origin 'origin contribution commit' }}:contribution |
      |              | git checkout observed                                                                            |
      | observed     | git reset --hard {{ sha 'local observed commit' }}                                               |
      |              | git checkout alpha                                                                               |
    And the initial branches and lineage exist now
    And the initial commits exist now
