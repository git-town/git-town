Feature: git-town sync --all: syncs all perennial branches

  Background:
    Given my repository has the perennial branches "production" and "qa"
    And the following commits exist in my repository
      | BRANCH     | LOCATION | MESSAGE                  | FILE NAME              |
      | main       | remote   | main commit              | main_file              |
      | production | local    | production local commit  | production_local_file  |
      |            | remote   | production remote commit | production_remote_file |
      | qa         | local    | qa local commit          | qa_local_file          |
      |            | remote   | qa remote commit         | qa_remote_file         |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town sync --all`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune            |
      |            | git add -A                   |
      |            | git stash                    |
      |            | git rebase origin/main       |
      |            | git checkout production      |
      | production | git rebase origin/production |
      |            | git push                     |
      |            | git checkout qa              |
      | qa         | git rebase origin/qa         |
      |            | git push                     |
      |            | git checkout main            |
      | main       | git push --tags              |
      |            | git stash pop                |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE                  | FILE NAME              |
      | main       | local and remote | main commit              | main_file              |
      | production | local and remote | production remote commit | production_remote_file |
      |            |                  | production local commit  | production_local_file  |
      | qa         | local and remote | qa remote commit         | qa_remote_file         |
      |            |                  | qa local commit          | qa_local_file          |
