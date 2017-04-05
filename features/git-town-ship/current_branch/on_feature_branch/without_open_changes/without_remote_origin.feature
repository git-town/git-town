Feature: git town-ship: shipping the current feature branch without a remote origin

  As a developer having finished a feature and on repo without a remote origin
  I want to be able to ship it safely in one easy step
  So that I can quickly move on to the next feature and remain productive.


  Background:
    Given I have a feature branch named "feature"
    And my repo does not have a remote origin
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run `gt ship -m "feature done"`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git merge --no-edit main     |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git branch -D feature        |
    And I end up on the "main" branch
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME    |
      | main   | local    | feature done | feature_file |


  Scenario: undo
    When I run `gt ship --undo`
    Then it runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git branch feature <%= sha 'feature commit' %> |
      |        | git revert HEAD           |
      |        | git checkout feature                           |
    And I end up on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME    |
      | main    | local    | feature done          | feature_file |
      |         |          | Revert "feature done" | feature_file |
      | feature | local    | feature commit        | feature_file |
