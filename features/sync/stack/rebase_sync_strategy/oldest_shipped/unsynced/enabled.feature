Feature: auto-resolve phantom merge conflicts after the oldest branch ships in an unsynced stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE                     | FILE NAME        | FILE CONTENT     |
      | branch-1 | local, origin | conflicting branch-1 commit | conflicting_file | branch-1 content |
      | branch-2 | local         | conflicting branch-2 commit | conflicting_file | branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
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
      |          | git checkout --theirs conflicting_file                     |
      |          | git add conflicting_file                                   |
      |          | GIT_EDITOR=true git rebase --continue                      |
      |          | git push --force-with-lease                                |
      |          | git branch -D branch-1                                     |
    And no rebase is now in progress

  Scenario: undo
    When I run "git town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                |
      | branch-2 | git reset --hard {{ sha-initial 'conflicting branch-2 commit' }}       |
      |          | git push --force-with-lease origin {{ sha 'initial commit' }}:branch-2 |
      |          | git checkout main                                                      |
      | main     | git reset --hard {{ sha 'initial commit' }}                            |
      |          | git branch branch-1 {{ sha-initial 'conflicting branch-1 commit' }}    |
      |          | git checkout branch-2                                                  |
    And the initial branches and lineage exist now
