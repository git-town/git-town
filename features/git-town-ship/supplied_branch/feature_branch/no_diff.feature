Feature: git town-ship: errors when trying to ship the supplied feature branch that has no differences with the main branch

  (see ../../current_branch/on_feature_branch/without_open_changes/no_diff.feature)


  Background:
    Given my repository has the feature branches "empty-feature" and "other-feature"
    And the following commit exists in my repository
      | BRANCH        | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main          | remote   | main commit    | common_file | common content |
      | empty-feature | local    | feature commit | common_file | common content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file
    When I run `git-town ship empty-feature`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH        | COMMAND                                      |
      | other-feature | git fetch --prune                            |
      |               | git add -A                                   |
      |               | git stash                                    |
      |               | git checkout main                            |
      | main          | git rebase origin/main                       |
      |               | git checkout empty-feature                   |
      | empty-feature | git merge --no-edit origin/empty-feature     |
      |               | git merge --no-edit main                     |
      |               | git reset --hard <%= sha 'feature commit' %> |
      |               | git checkout main                            |
      | main          | git checkout other-feature                   |
      | other-feature | git stash pop                                |
    And it prints the error "The branch 'empty-feature' has no shippable changes"
    And I am still on the "other-feature" branch
    And my workspace still contains my uncommitted file
