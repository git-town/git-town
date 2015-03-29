Feature: git sync --all: syncs all feature branches (without open changes or remote repo)

  Background:
    Given my repo does not have a remote origin
    And I have local feature branches named "feature1" and "feature2"
    And the following commits exist in my repository
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME     | FILE CONTENT     |
      | main     | local    | main commit     | main_file     | main content     |
      | feature1 | local    | feature1 commit | feature1_file | feature1 content |
      | feature2 | local    | feature2 commit | feature2_file | feature2 content |
    And I am on the "feature1" branch
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                  |
      | feature1 | git merge --no-edit main |
      |          | git checkout feature2    |
      | feature2 | git merge --no-edit main |
      |          | git checkout feature1    |
    And I am still on the "feature1" branch
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE                           | FILE NAME     |
      | main     | local    | main commit                       | main_file     |
      | feature1 | local    | feature1 commit                   | feature1_file |
      |          |          | main commit                       | main_file     |
      |          |          | Merge branch 'main' into feature1 |               |
      | feature2 | local    | feature2 commit                   | feature2_file |
      |          |          | main commit                       | main_file     |
      |          |          | Merge branch 'main' into feature2 |               |
    And now I have the following committed files
      | BRANCH   | NAME          | CONTENT          |
      | main     | main_file     | main content     |
      | feature1 | feature1_file | feature1 content |
      | feature1 | main_file     | main content     |
      | feature2 | feature2_file | feature2 content |
      | feature2 | main_file     | main content     |
