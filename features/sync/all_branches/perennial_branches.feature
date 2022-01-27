Feature: git-town sync --all: syncs all perennial branches

  Background:
    Given my repo has the perennial branches "production" and "qa"
    And the following commits exist in my repo
      | BRANCH     | LOCATION | MESSAGE                  |
      | main       | remote   | main commit              |
      | production | local    | production local commit  |
      |            | remote   | production remote commit |
      | qa         | local    | qa local commit          |
      |            | remote   | qa remote commit         |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune --tags     |
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
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | local, remote | main commit              |
      | production | local, remote | production remote commit |
      |            |               | production local commit  |
      | qa         | local, remote | qa remote commit         |
      |            |               | qa local commit          |
