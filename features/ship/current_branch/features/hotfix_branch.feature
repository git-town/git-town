Feature: ship hotfixes

  Background:
    Given my repo has the perennial branch "production"
    And my repo has a feature branch named "hotfix" as a child of "production"
    And the following commits exist in my repo
      | BRANCH | LOCATION      | MESSAGE       |
      | hotfix | local, remote | hotfix commit |
    And I am on the "hotfix" branch
    When I run "git-town ship -m 'hotfix done'"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                           |
      | hotfix     | git fetch --prune --tags          |
      |            | git checkout production           |
      | production | git rebase origin/production      |
      |            | git checkout hotfix               |
      | hotfix     | git merge --no-edit origin/hotfix |
      |            | git merge --no-edit production    |
      |            | git checkout production           |
      | production | git merge --squash hotfix         |
      |            | git commit -m "hotfix done"       |
      |            | git push                          |
      |            | git push origin :hotfix           |
      |            | git branch -D hotfix              |
    And I am now on the "production" branch
    And the existing branches are
      | REPOSITORY | BRANCHES         |
      | local      | main, production |
      | remote     | main, production |
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE     |
      | production | local, remote | hotfix done |
    And Git Town now has no branch hierarchy information

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                                     |
      | production | git branch hotfix {{ sha 'hotfix commit' }} |
      |            | git push -u origin hotfix                   |
      |            | git revert {{ sha 'hotfix done' }}          |
      |            | git push                                    |
      |            | git checkout hotfix                         |
      | hotfix     | git checkout production                     |
      | production | git checkout hotfix                         |
    And I am now on the "hotfix" branch
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE              |
      | hotfix     | local, remote | hotfix commit        |
      | production | local, remote | hotfix done          |
      |            |               | Revert "hotfix done" |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT     |
      | hotfix | production |
