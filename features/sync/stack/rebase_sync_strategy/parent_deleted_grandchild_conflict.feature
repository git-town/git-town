Feature: a grandchild branch has conflicts while its parent was deleted remotely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE    | PARENT | LOCATIONS     |
      | child      | feature | main   | local, origin |
      | grandchild | feature | child  | local, origin |
    And the commits
      | BRANCH     | LOCATION | MESSAGE                       | FILE NAME        | FILE CONTENT       |
      | main       | local    | conflicting main commit       | conflicting_file | main content       |
      | child      | local    | child commit                  | child_file       | child content      |
      | grandchild | local    | conflicting grandchild commit | conflicting_file | grandchild content |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And origin deletes the "child" branch
    And the current branch is "child" and the previous branch is "grandchild"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                   |
      | child      | git fetch --prune --tags                  |
      |            | git checkout main                         |
      | main       | git rebase origin/main --no-update-refs   |
      |            | git push                                  |
      |            | git checkout grandchild                   |
      | grandchild | git pull                                  |
      |            | git rebase --onto main child              |
      |            | git checkout --theirs conflicting_file    |
      |            | git add conflicting_file                  |
      |            | git -c core.editor=true rebase --continue |
      |            | git push --force-with-lease               |
      |            | git branch -D child                       |
      |            | git push --tags                           |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                       | FILE NAME        | FILE CONTENT       |
      | main       | local, origin | conflicting main commit       | conflicting_file | main content       |
      | grandchild | local, origin | conflicting grandchild commit | conflicting_file | grandchild content |
