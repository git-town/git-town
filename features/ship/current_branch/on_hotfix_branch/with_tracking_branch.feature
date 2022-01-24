Feature: git town-ship: shipping hotfixes


  Background:
    Given my repo has the perennial branch "production"
    And my repo has a feature branch named "hotfix" as a child of "production"
    And the following commits exist in my repo
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   | FILE CONTENT   |
      | hotfix | local, remote | hotfix commit | hotfix_file | hotfix content |
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
      | BRANCH     | LOCATION      | MESSAGE     | FILE NAME   |
      | production | local, remote | hotfix done | hotfix_file |
