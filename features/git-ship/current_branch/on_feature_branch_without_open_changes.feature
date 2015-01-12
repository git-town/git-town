Feature: git ship: shipping the current feature branch

  As a developer having finished a feature
  I want to be able to ship it safely in one easy step
  So that I can quickly move on to the next feature and remain productive.


  Scenario: local feature branch
    Given I have a local feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature_file | feature content |
    And I am on the "feature" branch
    When I run `git ship -m 'feature done'`
    Then it runs the Git commands
      | BRANCH  | COMMAND                      |
      | feature | git checkout main            |
      | main    | git fetch --prune            |
      | main    | git rebase origin/main       |
      | main    | git checkout feature         |
      | feature | git merge --no-edit main     |
      | feature | git checkout main            |
      | main    | git merge --squash feature   |
      | main    | git commit -m 'feature done' |
      | main    | git push                     |
      | main    | git branch -D feature        |
    And I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |


  Scenario: feature branch with non-pulled updates in the repo
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | FILE NAME    | FILE CONTENT          |
      | feature | local and remote | feature_file | early feature content |
      | feature | local and remote | feature_file | mid feature content   |
      | feature | remote           | feature_file | final feature content |
    And I am on the "feature" branch
    When I run `git ship -m 'feature done'`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git fetch --prune                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git checkout main                  |
      | main    | git merge --squash feature         |
      | main    | git commit -m 'feature done'       |
      | main    | git push                           |
      | main    | git push origin :feature           |
      | main    | git branch -D feature              |
    And I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
