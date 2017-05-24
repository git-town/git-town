Feature: git town-ship: shipping a coworker's feature branch

  As a developer shipping a coworker's feature branch
  I want my coworker to be the author of the commit added to the main branch
  So that my coworker is given credit for their work


  Background:
    Given my coworker has a feature branch named "feature"
    And the following commits exist in my coworker's repository
      | BRANCH  | LOCATION         | MESSAGE         | FILE NAME     |
      | feature | local and remote | coworker commit | coworker_file |
    And I fetch updates
    And I set the parent branch of "feature" as "main"
    And I am on the "feature" branch
    When I run `git-town ship -m 'feature done'`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                                                 |
      | feature | git fetch --prune                                                       |
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
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME     | AUTHOR                          |
      | main   | local and remote | feature done | coworker_file | coworker <coworker@example.com> |
