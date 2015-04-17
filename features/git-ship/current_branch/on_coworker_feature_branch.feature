Feature: git ship: shipping a coworker's feature branch

  As a developer shipping a coworker's feature branch
  I want my coworker to be the author of the commit added to the main branch
  So that my coworker is given credit for their work


  Background:
    Given my coworker has a feature branch named "feature"
    And the following commits exist in my coworker's repository
      | BRANCH  | LOCATION         | MESSAGE         | FILE NAME     |
      | feature | local and remote | coworker commit | coworker_file |
    And I fetch updates
    And I am on the "feature" branch
    When I run `git ship -m 'feature done'`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                                                                 |
      | feature | git checkout main                                                       |
      | main    | git fetch --prune                                                       |
      |         | git rebase origin/main                                                  |
      |         | git checkout feature                                                    |
      | feature | git merge --no-edit origin/feature                                      |
      |         | git merge --no-edit main                                                |
      |         | git checkout main                                                       |
      | main    | git merge --squash feature                                              |
      |         | git commit --author="coworker <coworker@primary.com>" -m "feature done" |
      |         | git push                                                                |
      |         | git push origin :feature                                                |
      |         | git branch -D feature                                                   |
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME     | AUTHOR               |
      | main   | local and remote | feature done | coworker_file | coworker@primary.com |
