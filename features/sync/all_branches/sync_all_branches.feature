Feature: sync all feature branches

  Background:
    Given the feature branches "alpha" and "beta"
    And the perennial branches "production" and "qa"
    And the commits
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | origin        | main commit              |
      | alpha      | local, origin | alpha commit             |
      | beta       | local, origin | beta commit              |
      | production | local         | local production commit  |
      |            | origin        | origin production commit |
      | qa         | local         | qa local commit          |
      |            | origin        | qa origin commit         |
    And the current branch is "alpha"
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                          |
      | alpha      | git fetch --prune --tags         |
      |            | git checkout main                |
      | main       | git rebase origin/main           |
      |            | git checkout alpha               |
      | alpha      | git merge --no-edit origin/alpha |
      |            | git merge --no-edit main         |
      |            | git push                         |
      |            | git checkout beta                |
      | beta       | git merge --no-edit origin/beta  |
      |            | git merge --no-edit main         |
      |            | git push                         |
      |            | git checkout production          |
      | production | git rebase origin/production     |
      |            | git push                         |
      |            | git checkout qa                  |
      | qa         | git rebase origin/qa             |
      |            | git push                         |
      |            | git checkout alpha               |
      | alpha      | git push --tags                  |
    And the current branch is still "alpha"
    And all branches are now synchronized
