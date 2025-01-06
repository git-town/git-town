Feature: do not undo branches that were created while resolving conflicts

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the current branch is "feature-1"
    And the commits
      | BRANCH    | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature-1 | local    | conflicting local commit  | conflicting_file | local content  |
      |           | origin   | conflicting origin commit | conflicting_file | origin content |
    And I run "git-town sync"
    And Git Town runs the commands
      | BRANCH    | COMMAND                                   |
      | feature-1 | git fetch --prune --tags                  |
      |           | git checkout main                         |
      | main      | git rebase origin/main --no-update-refs   |
      |           | git checkout feature-1                    |
      | feature-1 | git merge --no-edit --ff main             |
      |           | git merge --no-edit --ff origin/feature-1 |
    And Git Town prints the error:
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
    When I run "git-town continue"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND  |
      | feature-1 | git push |
    And no merge is in progress
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                                      |
      | feature-1 | git reset --hard {{ sha 'conflicting local commit' }}                                        |
      |           | git push --force-with-lease origin {{ sha-in-origin 'conflicting origin commit' }}:feature-1 |
    And no merge is in progress
    And the current branch is still "feature-1"
    And these branches exist now
      | REPOSITORY | BRANCHES                   |
      | local      | main, feature-1, feature-2 |
      | origin     | main, feature-1            |
    And these commits exist now
      | BRANCH    | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature-1 | local    | conflicting local commit  | conflicting_file | local content  |
      |           | origin   | conflicting origin commit | conflicting_file | origin content |
      | feature-2 | local    | feature-2 commit          | feature_2_file   | content 2      |
