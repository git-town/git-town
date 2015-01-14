Feature: git sync --all: syncs all non-feature branches without open changes

  Background:
    Given I have branches named "production" and "qa"
    And my non-feature branches are "production" and "qa"
    And the following commits exist in my repository
      | BRANCH     | LOCATION | MESSAGE                  | FILE NAME              |
      | main       | remote   | main commit              | main_file              |
      | production | local    | production local commit  | production_local_file  |
      |            | remote   | production remote commit | production_remote_file |
      | qa         | local    | qa local commit          | qa_local_file          |
      |            | remote   | qa remote commit         | qa_remote_file         |
    And I am on the "main" branch
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune            |
      | main       | git rebase origin/main       |
      | main       | git checkout production      |
      | production | git rebase origin/production |
      | production | git push                     |
      | production | git checkout qa              |
      | qa         | git rebase origin/qa         |
      | qa         | git push                     |
      | qa         | git checkout main            |
    And I am still on the "main" branch
    And all branches are now synchronized
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME              |
      | main       | local and remote | main commit              | main_file              |
      | production | local and remote | production remote commit | production_remote_file |
      |            | local and remote | production local commit  | production_local_file  |
      | qa         | local and remote | qa remote commit         | qa_remote_file         |
      |            | local and remote | qa local commit          | qa_local_file          |