Feature: handle a created branch while resolving conflicts

  Background: I fetch updates while resolving merge conflicts
    Given a Git repo clone
    And the branch
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the current branch is "feature-1"
    And the commits
      | BRANCH    | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature-1 | local    | conflicting local commit  | conflicting_file | local content  |
      |           | origin   | conflicting origin commit | conflicting_file | origin content |
    And I run "git-town sync"
    And it runs the commands
      | BRANCH    | COMMAND                                   |
      | feature-1 | git fetch --prune --tags                  |
      |           | git checkout main                         |
      | main      | git rebase origin/main                    |
      |           | git checkout feature-1                    |
      | feature-1 | git merge --no-edit --ff origin/feature-1 |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And the current branch is still "feature-1"
    And a merge is now in progress
    And I resolve the conflict in "conflicting_file"
    And I run "git add ."
    And I run "git commit --no-edit"
    And in a separate terminal I create branch "feature-2" with commits
      | MESSAGE          | FILE NAME      | FILE CONTENT |
      | feature-2 commit | feature_2_file | content 2    |
    When I run "git-town continue" and enter into the dialog:
      | DIALOG                      | KEYS  |
      | parent branch for feature-2 | enter |

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | feature-1 | git merge --no-edit --ff main |
      |           | git push                      |
    And no merge is in progress
    And all branches are now synchronized

  @this
  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                                                                                      |
      | feature-1 | git reset --hard {{ sha 'conflicting local commit' }}                                        |
      |           | git push --force-with-lease origin {{ sha-in-origin 'conflicting origin commit' }}:feature-1 |
      |           | git branch -D feature-2                                                                      |
    And no merge is in progress
    And the current branch is still "feature-1"
    And the initial branches and lineage exist
    And the initial commits exist
