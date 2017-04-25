Feature: git-town sync --all: handling merge conflicts between feature branch and main branch (without remote repo)

  Background:
    Given my repo does not have a remote origin
    And I have local feature branches named "feature-1" and "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION | MESSAGE          | FILE NAME        | FILE CONTENT      |
      | main      | local    | main commit      | conflicting_file | main content      |
      | feature-1 | local    | feature-1 commit | conflicting_file | feature-1 content |
      | feature-2 | local    | feature-2 commit | feature2_file    | feature-2 content |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git-town sync --all`


  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | main      | git add -A               |
      |           | git stash                |
      |           | git checkout feature-1   |
      | feature-1 | git merge --no-edit main |
    And I get the error
      """
      To abort, run "git-town sync --abort".
      To continue after you have resolved the conflicts, run "git-town sync --continue".
      To skip the sync of the 'feature-1' branch, run "git-town sync --skip".
      """
    And I end up on the "feature-1" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git-town sync --abort`
    Then it runs the commands
      | BRANCH    | COMMAND           |
      | feature-1 | git merge --abort |
      |           | git checkout main |
      | main      | git stash pop     |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I am left with my original commits


  Scenario: skipping
    When I run `git-town sync --skip`
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git merge --abort        |
      |           | git checkout feature-2   |
      | feature-2 | git merge --no-edit main |
      |           | git checkout main        |
      | main      | git stash pop            |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION | MESSAGE                            | FILE NAME        |
      | main      | local    | main commit                        | conflicting_file |
      | feature-1 | local    | feature-1 commit                   | conflicting_file |
      | feature-2 | local    | feature-2 commit                   | feature2_file    |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-2 |                  |
  And now I have the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | feature-1 content |
      | feature-2 | conflicting_file | main content      |
      | feature-2 | feature2_file    | feature-2 content |


  Scenario: continuing without resolving the conflicts
    When I run `git-town sync --continue`
    Then it runs no commands
    And I get the error "You must resolve the conflicts before continuing"
    And I am still on the "feature-1" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress
    And now I have the following committed files
        | BRANCH    | NAME             | CONTENT           |
        | main      | conflicting_file | main content      |
        | feature-1 | conflicting_file | feature-1 content |
        | feature-2 | feature2_file    | feature-2 content |


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git-town sync --continue`
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git commit --no-edit     |
      |           | git checkout feature-2   |
      | feature-2 | git merge --no-edit main |
      |           | git checkout main        |
      | main      | git stash pop            |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION | MESSAGE                            | FILE NAME        |
      | main      | local    | main commit                        | conflicting_file |
      | feature-1 | local    | feature-1 commit                   | conflicting_file |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-1 |                  |
      | feature-2 | local    | feature-2 commit                   | feature2_file    |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-2 |                  |
    And now I have the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | resolved content  |
      | feature-2 | conflicting_file | main content      |
      | feature-2 | feature2_file    | feature-2 content |


  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit; git-town sync --continue`
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git checkout feature-2   |
      | feature-2 | git merge --no-edit main |
      |           | git checkout main        |
      | main      | git stash pop            |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION | MESSAGE                            | FILE NAME        |
      | main      | local    | main commit                        | conflicting_file |
      | feature-1 | local    | feature-1 commit                   | conflicting_file |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-1 |                  |
      | feature-2 | local    | feature-2 commit                   | feature2_file    |
      |           |          | main commit                        | conflicting_file |
      |           |          | Merge branch 'main' into feature-2 |                  |
    And now I have the following committed files
      | BRANCH    | NAME             | CONTENT           |
      | main      | conflicting_file | main content      |
      | feature-1 | conflicting_file | resolved content  |
      | feature-2 | conflicting_file | main content      |
      | feature-2 | feature2_file    | feature-2 content |
