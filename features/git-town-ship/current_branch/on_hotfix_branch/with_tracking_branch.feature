Feature: git town-ship: shipping hotfixes

  When working on hotfix branches
  I want to ship them similar to feature branches
  So that I can use Git Town to work on hotfixes as well.


  Background:
    Given my repository has a perennial branch named "production"
    And my repository has a hotfix branch named "hotfix" as a child of "production"
    And the following commit exists in my repository
      | BRANCH | LOCATION         | MESSAGE       | FILE NAME   | FILE CONTENT   |
      | hotfix | local and remote | hotfix commit | hotfix_file | hotfix content |
    And I am on the "hotfix" branch
    When I run `git-town ship -m "hotfix done"`


  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                           |
      | hotfix     | git fetch --prune                 |
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
    And there are no more feature branches
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE     | FILE NAME   |
      | production | local and remote | hotfix done | hotfix_file |
