Feature: git extract: extracting multiple commits (with open changes)

  As a developer working on a feature branch with many commits around an unrelated issue
  I want to be able to extract all of these commits into their own branch
  So that the issue can be reviewed separately and my feature branch remains focussed.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor1 commit   | refactor1_file   |
      |         |          | refactor2 commit   | refactor2_file   |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` with the last two commit shas


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                                                       |
      | feature  | git fetch --prune                                             |
      | feature  | git stash -u                                                  |
      | feature  | git checkout main                                             |
      | main     | git rebase origin/main                                        |
      | main     | git checkout -b refactor main                                 |
      | refactor | git cherry-pick [SHA:refactor2 commit] [SHA:refactor1 commit] |
      | refactor | git push -u origin refactor                                   |
      | refactor | git stash pop                                                 |
    And I end up on the "refactor" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            | FILE NAME        |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
      | refactor | local and remote | remote main commit | remote_main_file |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
