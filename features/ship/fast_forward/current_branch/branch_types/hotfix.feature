Feature: ship hotfixes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | PARENT     | LOCATIONS     |
      | production | perennial |            | local, origin |
      | hotfix     | feature   | production | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | hotfix | local, origin | hotfix commit |
    And the current branch is "hotfix"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                    |
      | hotfix     | git fetch --prune --tags   |
      |            | git checkout production    |
      | production | git merge --ff-only hotfix |
      |            | git push                   |
      |            | git push origin :hotfix    |
      |            | git branch -D hotfix       |
    And the current branch is now "production"
    And the branches are now
      | REPOSITORY    | BRANCHES         |
      | local, origin | main, production |
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE       |
      | production | local, origin | hotfix commit |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                     |
      | production | git branch hotfix {{ sha 'hotfix commit' }} |
      |            | git push -u origin hotfix                   |
      |            | git checkout hotfix                         |
    And the current branch is now "hotfix"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE       |
      | production | local, origin | hotfix commit |
    And the initial branches and lineage exist now
