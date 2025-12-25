Feature: auto-resolve phantom merge conflicts in a synced stack where the parent deletes a file and gets shipped, and the child modifies the same file

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main   | local, origin | main commit | file      | main content |
    And I ran "git-town hack branch-1"
    And I ran "git rm file"
    And I ran "git commit -m branch-1-commit"
    And I ran "git push -u origin branch-1"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME | FILE CONTENT     |
      | branch-2 | local    | branch-2 commit | file      | branch-2 content |
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
    And no rebase is now in progress
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                         |
      | branch-2 | git reset --hard {{ sha-initial 'branch-2 commit' }}                            |
      |          | git push --force-with-lease origin {{ sha-initial 'branch-1-commit' }}:branch-2 |
      |          | git checkout main                                                               |
      | main     | git reset --hard {{ sha 'main commit' }}                                        |
      |          | git branch branch-1 {{ sha-initial 'branch-1-commit' }}                         |
      |          | git checkout branch-2                                                           |
    And no merge is now in progress
    And the initial commits exist now
