Feature: git town-ship: shipping the current feature branch

  As a developer entering a commit message that contains a double quote
  I want it to still work as expected
  So shipping is a robust process.


  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local, remote | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run "git-town ship -m 'message containing "double quotes"'"


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                                |
      | feature | git fetch --prune --tags                               |
      |         | git checkout main                                      |
      | main    | git rebase origin/main                                 |
      |         | git checkout feature                                   |
      | feature | git merge --no-edit origin/feature                     |
      |         | git merge --no-edit main                               |
      |         | git checkout main                                      |
      | main    | git merge --squash feature                             |
      |         | git commit -m "message containing \\"double quotes\\"" |
      |         | git push                                               |
      |         | git push origin :feature                               |
      |         | git branch -D feature                                  |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |
    And I don't have any uncommitted files
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE                            | FILE NAME    |
      | main   | local, remote | message containing "double quotes" | feature_file |


  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                   |
      | main    | git branch feature {{ sha 'feature commit' }}             |
      |         | git push -u origin feature                                |
      |         | git revert {{ sha 'message containing "double quotes"' }} |
      |         | git push                                                  |
      |         | git checkout feature                                      |
      | feature | git checkout main                                         |
      | main    | git checkout feature                                      |
    And I end up on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                                     | FILE NAME    |
      | main    | local, remote | message containing "double quotes"          | feature_file |
      |         |               | Revert "message containing "double quotes"" | feature_file |
      | feature | local, remote | feature commit                              | feature_file |
