Feature: git ship: don't ship a feature branch without changes (without open changes)

  (see ../current_branch/no_diff.feature)


  Background:
    Given I have feature branches named "empty-feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH        | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main          | remote   | main commit    | common_file | common content |
      | empty-feature | local    | feature commit | common_file | common content |
    And I am on the "other_feature" branch
    When I run `git ship empty-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH        | COMMAND                                  |
      | other_feature | git checkout main                        |
      | main          | git fetch --prune                        |
      | main          | git rebase origin/main                   |
      | main          | git checkout empty-feature               |
      | empty-feature | git merge --no-edit origin/empty-feature |
      | empty-feature | git merge --no-edit main                 |
      | empty-feature | git reset --hard [SHA:feature commit]    |
      | empty-feature | git checkout main                        |
      | main          | git checkout other_feature               |
    And I get the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "other_feature" branch
