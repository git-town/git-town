Feature: sync the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE       | FILE NAME   |
      | main         | local, origin | main commit   | main_file   |
      | contribution | local         | local commit  | local_file  |
      |              | origin        | origin commit | origin_file |
    And the current branch is "contribution"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                   |
      | contribution | git fetch --prune --tags                                  |
      |              | git -c rebase.updateRefs=false rebase origin/contribution |
      |              | git push                                                  |
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE       |
      | main         | local, origin | main commit   |
      | contribution | local, origin | origin commit |
      |              |               | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                                             |
      | contribution | git reset --hard {{ sha-initial 'local commit' }}                                   |
      |              | git push --force-with-lease origin {{ sha-in-origin 'origin commit' }}:contribution |
    And the initial branches and lineage exist now
    And the initial commits exist now
