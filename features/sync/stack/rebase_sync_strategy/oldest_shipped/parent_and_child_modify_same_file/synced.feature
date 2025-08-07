Feature: sync a synced stack where parent modifies a file and gets shipped, and the child modifies the same file

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT       |
      | main   | local, origin | main commit | file      | line 1 \n\n line 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                           |
      | branch-1 | local, origin | branch-1 commit | file      | line 1 changed by branch-1 \n\n line 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | branch-2 | local    | branch-2 commit | file      | line 1 changed by branch-1 \n\n line 2 changed by branch-2 |
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git checkout main                                          |
      | main     | git -c rebase.updateRefs=false rebase origin/main          |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git push --force-with-lease                                |
      |          | git branch -D branch-1                                     |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | main     | local, origin | main commit     | file      | line 1 \n\n line 2                                         |
      |          |               | branch-1 commit | file      | line 1 changed by branch-1 \n\n line 2                     |
      | branch-2 | local, origin | branch-2 commit | file      | line 1 changed by branch-1 \n\n line 2 changed by branch-2 |
    And no rebase is now in progress

  Scenario: undo
    When I run "git town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                         |
      | branch-2 | git reset --hard {{ sha-initial 'branch-2 commit' }}                            |
      |          | git push --force-with-lease origin {{ sha-initial 'branch-1 commit' }}:branch-2 |
      |          | git checkout main                                                               |
      | main     | git reset --hard {{ sha 'main commit' }}                                        |
      |          | git branch branch-1 {{ sha-initial 'branch-1 commit' }}                         |
      |          | git checkout branch-2                                                           |
    And the initial branches and lineage exist now
