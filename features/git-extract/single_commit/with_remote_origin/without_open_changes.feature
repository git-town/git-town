Feature: git extract: extracting a single commit (without open changes)

  (see ../../multiple_commits/with_remote_origin/with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    When I run `git extract refactor` with the last commit sha


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                                      |
      | feature  | git fetch --prune                            |
      |          | git checkout main                            |
      | main     | git rebase origin/main                       |
      |          | git checkout -b refactor main                |
      | refactor | git cherry-pick <%= sha 'refactor commit' %> |
      |          | git push -u origin refactor                  |
    And I end up on the "refactor" branch
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            |
      | main     | local and remote | remote main commit |
      | feature  | local            | feature commit     |
      |          |                  | refactor commit    |
      | refactor | local and remote | remote main commit |
      |          |                  | refactor commit    |
    And now I have the following committed files
      | BRANCH   | NAME             |
      | main     | remote_main_file |
      | feature  | feature_file     |
      | feature  | refactor_file    |
      | refactor | refactor_file    |
      | refactor | remote_main_file |
