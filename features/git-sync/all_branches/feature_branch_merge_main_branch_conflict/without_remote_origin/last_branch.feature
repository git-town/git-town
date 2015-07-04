Feature: git sync --all: handling merge conflicts between feature branch and main branch (without remote repo)

  Background:
    Given my repo does not have a remote origin
    And I have local feature branches named "feature-1" and "feature2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION | MESSAGE          | FILE NAME        | FILE CONTENT      |
      | main      | local    | main commit      | conflicting_file | main content      |
      | feature-1 | local    | feature-1 commit | feature1_file    | feature-1 content |
      | feature2  | local    | feature2 commit  | conflicting_file | feature2 content  |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH    | COMMAND                  |
      | main      | git stash -u             |
      |           | git checkout feature-1   |
      | feature-1 | git merge --no-edit main |
      |           | git checkout feature2    |
      | feature2  | git merge --no-edit main |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      To skip the sync of the 'feature2' branch, run "git sync --skip".
      """
    And I end up on the "feature2" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH    | COMMAND                                        |
      | feature2  | git merge --abort                              |
      |           | git checkout feature-1                         |
      | feature-1 | git reset --hard <%= sha 'feature-1 commit' %> |
      |           | git checkout main                              |
      | main      | git stash pop                                  |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I am left with my original commits


  Scenario: skipping
    When I run `git sync --skip`
    Then it runs the Git commands
      | BRANCH   | COMMAND           |
      | feature2 | git merge --abort |
      |          | git checkout main |
      | main     | git stash pop     |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION | MESSAGE                            | FILE NAME        |
      | main      | local    | main commit                        | conflicting_file |
      | feature-1 | local    | feature-1 commit                   | feature1_file    |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-1 |                  |
      | feature2  | local    | feature2 commit                    | conflicting_file |
    And now I have the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | main content      |
      | feature-1 | feature1_file    | feature-1 content |
      | feature2  | conflicting_file | feature2 content  |


  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature2" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND              |
      | feature2 | git commit --no-edit |
      |          | git checkout main    |
      | main     | git stash pop        |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION | MESSAGE                            | FILE NAME        |
      | main      | local    | main commit                        | conflicting_file |
      | feature-1 | local    | feature-1 commit                   | feature1_file    |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-1 |                  |
      | feature2  | local    | feature2 commit                    | conflicting_file |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature2  |                  |
    And now I have the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | main content      |
      | feature-1 | feature1_file    | feature-1 content |
      | feature2  | conflicting_file | resolved content  |



  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit; git sync --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND           |
      | feature2 | git checkout main |
      | main     | git stash pop     |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION | MESSAGE                            | FILE NAME        |
      | main      | local    | main commit                        | conflicting_file |
      | feature-1 | local    | feature-1 commit                   | feature1_file    |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-1 |                  |
      | feature2  | local    | feature2 commit                    | conflicting_file |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature2  |                  |
    And now I have the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | main content      |
      | feature-1 | feature1_file    | feature-1 content |
      | feature2  | conflicting_file | resolved content  |
