Feature: git sync --all: syncs all feature branches (with open changes and without remote repo)

  Background:
    Given I have feature branches named "feature1" and "feature2"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME     |
      | main     | local    | main commit     | main_file     |
      | feature1 | local    | feature1 commit | feature1_file |
      | feature2 | local    | feature2 commit | feature2_file |
    And I am on the "feature1" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                  |
      | feature1 | git stash -u             |
      | feature1 | git merge --no-edit main |
      | feature1 | git checkout feature2    |
      | feature2 | git merge --no-edit main |
      | feature2 | git checkout feature1    |
      | feature1 | git stash pop            |
    And I am still on the "feature1" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE                           | FILE NAME     |
      | main     | local    | main commit                       | main_file     |
      | feature1 | local    | feature1 commit                   | feature1_file |
      |          |          | main commit                       | main_file     |
      |          |          | Merge branch 'main' into feature1 |               |
      | feature2 | local    | feature2 commit                   | feature2_file |
      |          |          | main commit                       | main_file     |
      |          |          | Merge branch 'main' into feature2 |               |
