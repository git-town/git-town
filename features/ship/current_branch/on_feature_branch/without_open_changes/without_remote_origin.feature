Feature: git town-ship: shipping the current feature branch without a remote origin


  Background:
    Given my repo has a feature branch named "feature"
    And my repo does not have a remote origin
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run "git-town ship -m 'feature done'"


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git merge --no-edit main     |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git branch -D feature        |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
    And my repo now has the following commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME    |
      | main   | local    | feature done | feature_file |


  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
      |        | git revert {{ sha 'feature done' }}           |
      |        | git checkout feature                          |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME    |
      | main    | local    | feature done          | feature_file |
      |         |          | Revert "feature done" | feature_file |
      | feature | local    | feature commit        | feature_file |
