Feature: git town-ship: shipping a coworker's feature branch


  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE         | FILE NAME     | AUTHOR                          |
      | feature | local, remote | coworker commit | coworker_file | coworker <coworker@example.com> |
    And I am on the "feature" branch

  Scenario: result (commit message via CLI)
    When I run "git-town ship -m 'feature done'"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                 |
      | feature | git fetch --prune --tags                                                |
      |         | git checkout main                                                       |
      | main    | git rebase origin/main                                                  |
      |         | git checkout feature                                                    |
      | feature | git merge --no-edit origin/feature                                      |
      |         | git merge --no-edit main                                                |
      |         | git checkout main                                                       |
      | main    | git merge --squash feature                                              |
      |         | git commit -m "feature done" --author "coworker <coworker@example.com>" |
      |         | git push                                                                |
      |         | git push origin :feature                                                |
      |         | git branch -D feature                                                   |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME     | AUTHOR                          |
      | main   | local, remote | feature done | coworker_file | coworker <coworker@example.com> |

  Scenario: result (commit message via editor)
    When I run "git-town ship" and enter "feature done" for the commit message
    Then it runs the commands
      | BRANCH  | COMMAND                                               |
      | feature | git fetch --prune --tags                              |
      |         | git checkout main                                     |
      | main    | git rebase origin/main                                |
      |         | git checkout feature                                  |
      | feature | git merge --no-edit origin/feature                    |
      |         | git merge --no-edit main                              |
      |         | git checkout main                                     |
      | main    | git merge --squash feature                            |
      |         | git commit --author "coworker <coworker@example.com>" |
      |         | git push                                              |
      |         | git push origin :feature                              |
      |         | git branch -D feature                                 |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME     | AUTHOR                          |
      | main   | local, remote | feature done | coworker_file | coworker <coworker@example.com> |
