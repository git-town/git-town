Feature: git town-ship: shipping a coworker's feature branch

  As a developer shipping a coworker's feature branch
  I want my coworker to be the author of the commit added to the main branch
  So that my coworker is given credit for their work


  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE         | FILE NAME     | AUTHOR                          |
      | feature | local, remote | coworker commit | coworker_file | coworker <coworker@example.com> |
    And I am on the "feature" branch
    When I run "git-town ship -m 'feature done'"


  Scenario: result
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
