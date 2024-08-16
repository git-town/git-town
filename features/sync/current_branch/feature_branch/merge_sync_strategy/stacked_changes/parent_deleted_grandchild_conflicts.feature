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
      | grandchild | local    | conflicting grandchild commit | conflicting_file | grandchild content |
    And origin deletes the "child" branch
    And the current branch is "child" and the previous branch is "grandchild"
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                                    |
      | child      | git fetch --prune --tags                   |
      |            | git checkout main                          |
      | main       | git rebase origin/main                     |
      |            | git push                                   |
      |            | git checkout child                         |
      | child      | git merge --no-edit --ff main              |
      |            | git checkout main                          |
      | main       | git branch -D child                        |
      |            | git checkout grandchild                    |
      | grandchild | git merge --no-edit --ff origin/grandchild |
      |            | git merge --no-edit --ff main              |
    And it prints the error:
      """
      exit status 1
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git town skip".
      """
    And the current branch is now "grandchild"
    And a merge is now in progress

  Scenario: skip the grandchild merge conflict and kill the grandchild branch
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH     | COMMAND           |
      | grandchild | git merge --abort |
      |            | git push --tags   |
    And the current branch is now "grandchild"
    When I run "git-town kill"
    Then it runs the commands
      | BRANCH     | COMMAND                     |
      | grandchild | git fetch --prune --tags    |
      |            | git push origin :grandchild |
      |            | git checkout main           |
      | main       | git branch -D grandchild    |
