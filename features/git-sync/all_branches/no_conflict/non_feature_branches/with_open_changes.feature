Feature: git sync --all: syncs all non-feature branches with open changes

  Background:
    Given I have branches named "production" and "qa"
    And my non-feature branches are configured as "production" and "qa"
    And the following commits exist in my repository
      | BRANCH     | LOCATION | MESSAGE                  | FILE NAME              |
      | main       | remote   | main commit              | main_file              |
      | production | local    | production local commit  | production_local_file  |
      |            | remote   | production remote commit | production_remote_file |
      | qa         | local    | qa local commit          | qa_local_file          |
      |            | remote   | qa remote commit         | qa_remote_file         |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune            |
      |            | git stash -u                 |
      |            | git rebase origin/main       |
      |            | git checkout production      |
      | production | git rebase origin/production |
      |            | git push                     |
      |            | git checkout qa              |
      | qa         | git rebase origin/qa         |
      |            | git push                     |
      |            | git checkout main            |
      | main       | git stash pop                |
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And all branches are now synchronized
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME              |
      | main       | local and remote | main commit              | main_file              |
      | production | local and remote | production remote commit | production_remote_file |
      |            |                  | production local commit  | production_local_file  |
      | qa         | local and remote | qa remote commit         | qa_remote_file         |
      |            |                  | qa local commit          | qa_local_file          |
