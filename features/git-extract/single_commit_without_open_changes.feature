Feature: git extract: extracting a single commit (without open changes)

  (see ./multiple_commits_with_open_changes.feature)


  Background:
    Given I am on a feature branch
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      |         |          | refactor commit    | refactor_file    |
    When I run `git extract refactor` with the last commit sha


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                               |
      | feature  | git fetch --prune                     |
      | feature  | git checkout main                     |
      | main     | git rebase origin/main                |
      | main     | git checkout -b refactor main         |
      | refactor | git cherry-pick [SHA:refactor commit] |
      | refactor | git push -u origin refactor           |
    And I end up on the "refactor" branch
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE            | FILES            |
      | main     | local and remote | remote main commit | remote_main_file |
      | feature  | local            | feature commit     | feature_file     |
      |          |                  | refactor commit    | refactor_file    |
      | refactor | local and remote | remote main commit | remote_main_file |
      |          |                  | refactor commit    | refactor_file    |
    And now I have the following committed files
      | BRANCH   | FILES                           |
      | main     | remote_main_file                |
      | feature  | feature_file, refactor_file     |
      | refactor | remote_main_file, refactor_file |
