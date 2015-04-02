Feature: git extract: extracting multiple commits (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor1 commit   | refactor1_file   |
      |         |          | refactor2 commit   | refactor2_file   |
    And I am on the "feature" branch
    When I run `git extract refactor` with the last two commit shas


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                                                                     |
      | feature  | git fetch --prune                                                           |
      |          | git checkout main                                                           |
      | main     | git rebase origin/main                                                      |
      |          | git checkout -b refactor main                                               |
      | refactor | git cherry-pick <%= sha 'refactor1 commit' %> <%= sha 'refactor2 commit' %> |
      |          | git push -u origin refactor                                                 |
    And  I end up on the "refactor" branch
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            | FILE NAME        |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
      | refactor | local and remote | remote main commit | remote_main_file |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
    And now I have the following committed files
      | BRANCH   | NAME             |
      | main     | remote_main_file |
      | feature  | feature_file     |
      | feature  | refactor1_file   |
      | feature  | refactor2_file   |
      | refactor | refactor1_file   |
      | refactor | refactor2_file   |
      | refactor | remote_main_file |
