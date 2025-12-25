Feature: an ancestor in a stack with dependent changes was deleted remotely

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT       |
      | main   | local, origin | main commit | file      | line 1 \n\n line 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                         |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes \n\n line 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME | FILE CONTENT                                           |
      | branch-2 | local    | branch-2 commit | file      | line 1: branch-1 changes \n\n line 2: branch-2 changes |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin deletes the "branch-1" branch
    And the current branch is "branch-1" and the previous branch is "branch-2"
    When I run "git-town sync --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-1 | git fetch --prune --tags                                   |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git push --force-with-lease                                |
      |          | git branch -D branch-1                                     |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                         |
      | main     | local, origin | main commit     | file      | line 1 \n\n line 2                   |
      | branch-2 | local, origin | branch-2 commit | file      | line 1 \n\n line 2: branch-2 changes |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                 |
      | branch-2 | git reset --hard {{ sha 'branch-2 commit' }}                            |
      |          | git push --force-with-lease origin {{ sha 'branch-1 commit' }}:branch-2 |
      |          | git branch branch-1 {{ sha 'branch-1 commit' }}                         |
      |          | git checkout branch-1                                                   |
    And the initial branches and lineage exist now
