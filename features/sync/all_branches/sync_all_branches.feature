Feature: sync all feature branches

  Background:
    Given my repo has the feature branches "feature-1" and "feature-2"
    And my repo has the perennial branches "production" and "qa"
    And my repo contains the commits
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | remote        | main commit              |
      | feature-1  | local, remote | feature-1 commit         |
      | feature-2  | local, remote | feature-2 commit         |
      | production | local         | production local commit  |
      |            | remote        | production remote commit |
      | qa         | local         | qa local commit          |
      |            | remote        | qa remote commit         |
    And I am on the "feature-1" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                              |
      | feature-1  | git fetch --prune --tags             |
      |            | git add -A                           |
      |            | git stash                            |
      |            | git checkout main                    |
      | main       | git rebase origin/main               |
      |            | git checkout feature-1               |
      | feature-1  | git merge --no-edit origin/feature-1 |
      |            | git merge --no-edit main             |
      |            | git push                             |
      |            | git checkout feature-2               |
      | feature-2  | git merge --no-edit origin/feature-2 |
      |            | git merge --no-edit main             |
      |            | git push                             |
      |            | git checkout production              |
      | production | git rebase origin/production         |
      |            | git push                             |
      |            | git checkout qa                      |
      | qa         | git rebase origin/qa                 |
      |            | git push                             |
      |            | git checkout feature-1               |
      | feature-1  | git push --tags                      |
      |            | git stash pop                        |
    And I am still on the "feature-1" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
