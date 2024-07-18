Feature: ship hotfixes

  Background:
    Given a Git repo clone
    And the branches
      | NAME       | TYPE      | PARENT     | LOCATIONS     |
      | production | perennial |            | local, origin |
      | hotfix     | feature   | production | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | hotfix | local, origin | hotfix commit |
    And the current branch is "hotfix"
    When I run "git-town ship -m 'hotfix done'"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                        |
      | hotfix     | git fetch --prune --tags       |
      |            | git checkout production        |
      | production | git merge --squash --ff hotfix |
      |            | git commit -m "hotfix done"    |
      |            | git push                       |
      |            | git push origin :hotfix        |
      |            | git branch -D hotfix           |
    And the current branch is now "production"
    And the branches are now
      | REPOSITORY    | BRANCHES         |
      | local, origin | main, production |
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE     |
      | production | local, origin | hotfix done |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                                     |
      | production | git revert {{ sha 'hotfix done' }}          |
      |            | git push                                    |
      |            | git branch hotfix {{ sha 'hotfix commit' }} |
      |            | git push -u origin hotfix                   |
      |            | git checkout hotfix                         |
    And the current branch is now "hotfix"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE              |
      | hotfix     | local, origin | hotfix commit        |
      | production | local, origin | hotfix done          |
      |            |               | Revert "hotfix done" |
    And the initial branches and lineage exist
