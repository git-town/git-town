Feature: git extract: extracting multiple commits (without open changes)

  (see ./multiple_commits_with_open_changes.feature)


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
      | BRANCH   | COMMAND                                                       |
      | feature  | git fetch --prune                                             |
      | feature  | git checkout main                                             |
      | main     | git rebase origin/main                                        |
      | main     | git checkout -b refactor main                                 |
      | refactor | git cherry-pick [SHA:refactor2 commit] [SHA:refactor1 commit] |
      | refactor | git push -u origin refactor                                   |
    And  I end up on the "refactor" branch
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            | FILES            |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
      | refactor | local and remote | remote main commit | remote_main_file |
      |          |                  | refactor1 commit   | refactor1_file   |
      |          |                  | refactor2 commit   | refactor2_file   |
    And now I have the following committed files
      | BRANCH   | FILES                                            |
      | main     | remote_main_file                                 |
      | feature  | feature_file, refactor1_file, refactor2_file     |
      | refactor | remote_main_file, refactor1_file, refactor2_file |
