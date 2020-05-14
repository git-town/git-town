Feature: git town-ship: shipping hotfixes

  When working on hotfix branches
  I want to ship them similar to feature branches
  So that I can use Git Town to work on hotfixes as well.


  Background:
    Given my repository has the perennial branch "production"
    And my repository has a feature branch named "hotfix" as a child of "production"
    And the following commits exist in my repository
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
    And I end up on the "production" branch
    And the existing branches are
      | REPOSITORY | BRANCHES         |
      | local      | main, production |
      | remote     | main, production |
    And my repository now has the following commits
      | BRANCH     | LOCATION      | MESSAGE     | FILE NAME   |
      | production | local, remote | hotfix done | hotfix_file |
