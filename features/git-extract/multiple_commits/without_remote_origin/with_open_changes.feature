Feature: git extract: extracting multiple commits (with open changes and without remote repo)

  (see ../with_remote_origin/with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE          | FILE NAME      |
      | main    | local    | main commit      | main_file      |
      | feature | local    | feature commit   | feature_file   |
      |         |          | refactor1 commit | refactor1_file |
      |         |          | refactor2 commit | refactor2_file |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` with the last two commit shas


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                                                                     |
      | feature  | git stash -u                                                                |
      | feature  | git checkout main                                                           |
      | main     | git checkout -b refactor main                                               |
      | refactor | git cherry-pick <%= sha 'refactor1 commit' %> <%= sha 'refactor2 commit' %> |
      | refactor | git stash pop                                                               |
    And I end up on the "refactor" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE          | FILE NAME      |
      | main     | local    | main commit      | main_file      |
      | feature  | local    | feature commit   | feature_file   |
      |          |          | refactor1 commit | refactor1_file |
      |          |          | refactor2 commit | refactor2_file |
      | refactor | local    | main commit      | main_file      |
      |          |          | refactor1 commit | refactor1_file |
      |          |          | refactor2 commit | refactor2_file |
