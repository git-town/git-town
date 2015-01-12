Feature: git ship: shipping a coworkers branch

  As a developer shipping a coworker's feature branch
  I want my coworker to be the author of the commit added to the main branch
  So my coworker is given credit for their work


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
      | main    | git rebase origin/main                                                  |
      | main    | git checkout feature                                                    |
      | feature | git merge --no-edit origin/feature                                      |
      | feature | git merge --no-edit main                                                |
      | feature | git checkout main                                                       |
      | main    | git merge --squash feature                                              |
      | main    | git commit --author="coworker <coworker@example.com>" -m 'feature done' |
      | main    | git push                                                                |
      | main    | git push origin :feature                                                |
      | main    | git branch -D feature                                                   |
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME     | AUTHOR               |
      | main   | local and remote | feature done | coworker_file | coworker@example.com |
