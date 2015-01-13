Feature: git ship: shipping the supplied feature branch (without open changes)

  (see ./no_conflicts_with_conflicting_changes.feature)


  Scenario: local feature branch
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature_file | feature content |
    And I am on the "other_feature" branch
    When I run `git ship feature -m "feature done"`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      | feature       | git merge --no-edit main           |
      | feature       | git checkout main                  |
      | main          | git merge --squash feature         |
      | main          | git commit -m "feature done"       |
      | main          | git push                           |
      | main          | git push origin :feature           |
      | main          | git branch -D feature              |
      | main          | git checkout other_feature         |
    And I end up on the "other_feature" branch
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |


  Scenario: feature branch with non-pulled updates in the repo
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "other_feature" branch
    When I run `git ship feature -m "feature done"`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      | feature       | git merge --no-edit main           |
      | feature       | git checkout main                  |
      | main          | git merge --squash feature         |
      | main          | git commit -m "feature done"       |
      | main          | git push                           |
      | main          | git push origin :feature           |
      | main          | git branch -D feature              |
      | main          | git checkout other_feature         |
    And I end up on the "other_feature" branch
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
