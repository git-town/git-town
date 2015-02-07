Feature: git ship: shipping the current feature branch without a tracking branch

  (see ./with_tracking_branch.feature)


  Background:
    Given I have a local feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature_file | feature content |
    And I am on the "feature" branch
    When I run `git ship -m "feature done"`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                      |
      | feature | git checkout main            |
      | main    | git fetch --prune            |
      | main    | git rebase origin/main       |
      | main    | git checkout feature         |
      | feature | git merge --no-edit main     |
      | feature | git checkout main            |
      | main    | git merge --squash feature   |
      | main    | git commit -m "feature done" |
      | main    | git push                     |
      | main    | git branch -D feature        |
    And I end up on the "main" branch
    And there are no more feature branches
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
