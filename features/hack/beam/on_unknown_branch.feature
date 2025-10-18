@messyoutput
Feature: beam multiple commits from an unknown branch

  Background:
    Given a Git repo with origin
    And I ran "git checkout -b branch-1"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | branch-1 | local    | branch commit 1 |
      | branch-1 | local    | branch commit 2 |
    And the current branch is "branch-1"
    When I run "git-town hack --beam branch-2"
    # When I run "git-town hack --beam branch-2" and enter into the dialog:
    #   | DIALOG          | KEYS                   |
    #   | parent branch   | enter                  |
    #   | commits to beam | space down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                       |
      | branch-1 | git checkout -b branch-2 main |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        branch-2
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                |
      | new      | git checkout existing                                                  |
      | existing | git reset --hard {{ sha-initial 'commit 4' }}                          |
      |          | git push --force-with-lease origin {{ sha 'initial commit' }}:existing |
      |          | git branch -D new                                                      |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: amend the beamed commit
    And I amend this commit
      | BRANCH | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT    |
      | new    | local    | commit 4b | file_4    | amended content |
    And the current branch is "new"
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | new    | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout new                                  |
      | new    | git merge --no-edit --ff main                     |
      |        | git push -u origin new                            |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                      |
      | main     | local, origin | main commit                  |
      | existing | local, origin | commit 2                     |
      |          |               | commit 3                     |
      | new      | local, origin | commit 1                     |
      |          |               | commit 4b                    |
      |          |               | Merge branch 'main' into new |
