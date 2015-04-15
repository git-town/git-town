Feature: git ship: errors when trying to ship the supplied feature branch that has no differences with the main branch

  (see ../../current_branch/on_feature_branch/without_open_changes/no_diff.feature)


  Background:
    Given I have feature branches named "empty-feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH        | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main          | remote   | main commit    | common_file | common content |
      | empty-feature | local    | feature commit | common_file | common content |
    And I am on the "other_feature" branch


  Scenario: with open changes
    Given I have an uncommitted file
    When I run `git ship empty-feature`
    Then it runs the Git commands
      | BRANCH        | COMMAND                                      |
      | other_feature | git stash -u                                 |
      |               | git checkout main                            |
      | main          | git fetch --prune                            |
      |               | git rebase origin/main                       |
      |               | git checkout empty-feature                   |
      | empty-feature | git merge --no-edit origin/empty-feature     |
      |               | git merge --no-edit main                     |
      |               | git reset --hard <%= sha 'feature commit' %> |
      |               | git checkout main                            |
      | main          | git checkout other_feature                   |
      | other_feature | git stash pop                                |
    And I get the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "other_feature" branch
    And I still have my uncommitted file


  Scenario: without open changes
    When I run `git ship empty-feature`
    Then it runs the Git commands
      | BRANCH        | COMMAND                                      |
      | other_feature | git checkout main                            |
      | main          | git fetch --prune                            |
      |               | git rebase origin/main                       |
      |               | git checkout empty-feature                   |
      | empty-feature | git merge --no-edit origin/empty-feature     |
      |               | git merge --no-edit main                     |
      |               | git reset --hard <%= sha 'feature commit' %> |
      |               | git checkout main                            |
      | main          | git checkout other_feature                   |
    And I get the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "other_feature" branch
