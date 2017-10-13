Feature: git-town sync --all: syncs all feature branches (without remote repo)

  Background:
    Given my repo does not have a remote origin
    And my repository has the local feature branches "feature-1" and "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION | MESSAGE          | FILE NAME     | FILE CONTENT      |
      | main      | local    | main commit      | main_file     | main content      |
      | feature-1 | local    | feature-1 commit | feature1_file | feature-1 content |
      | feature-2 | local    | feature-2 commit | feature2_file | feature-2 content |
    And I am on the "feature-1" branch
    And my workspace has an uncommitted file
    When I run `git-town sync --all`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git add -A               |
      |           | git stash                |
      |           | git merge --no-edit main |
      |           | git checkout feature-2   |
      | feature-2 | git merge --no-edit main |
      |           | git checkout feature-1   |
      | feature-1 | git stash pop            |
    And I am still on the "feature-1" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH    | LOCATION | MESSAGE                            | FILE NAME     |
      | main      | local    | main commit                        | main_file     |
      | feature-1 | local    | feature-1 commit                   | feature1_file |
      |           |          | main commit                        | main_file     |
      |           |          | Merge branch 'main' into feature-1 |               |
      | feature-2 | local    | feature-2 commit                   | feature2_file |
      |           |          | main commit                        | main_file     |
      |           |          | Merge branch 'main' into feature-2 |               |
    And now my repository has the following committed files
      | BRANCH    | NAME          | CONTENT           |
      | main      | main_file     | main content      |
      | feature-1 | feature1_file | feature-1 content |
      | feature-1 | main_file     | main content      |
      | feature-2 | feature2_file | feature-2 content |
      | feature-2 | main_file     | main content      |
