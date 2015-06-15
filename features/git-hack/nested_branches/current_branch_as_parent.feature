Feature: Creating nested feature branches

  As a developer waiting for permission to ship a feature branch that contains changes needed for the next feature
  I want to be able to start working on the next feature while having access to the changes currently under review
  So that I am not slowed down by reviews and can keep working on my backlog.


  Background:
    Given I have a feature branch named "parent-feature"
    And Git Town is aware of this branch hierarchy
      | BRANCH         | PARENT |
      | parent-feature | main   |
    Given the following commits exist in my repository
      | BRANCH         | LOCATION | MESSAGE        | FILE NAME    |
      | main           | remote   | main_commit    | main_file    |
      | parent-feature | local    | feature_commit | feature_file |
    And I am on the "parent-feature" branch
    And I have an uncommitted file


  Scenario: Providing the name of the current branch
    When I run `git hack child-feature parent-feature`
    Then it runs the Git commands
      | BRANCH         | COMMAND                                      |
      | parent-feature | git fetch --prune                            |
      |                | git stash -u                                 |
      |                | git checkout main                            |
      | main           | git rebase origin/main                       |
      |                | git checkout parent-feature                  |
      | parent-feature | git merge --no-edit origin/parent-feature    |
      |                | git merge --no-edit main                     |
      |                | git push                                     |
      |                | git checkout -b child-feature parent-feature |
      | child-feature  | git stash pop                                |
    And I end up on the "child-feature" branch
    And I still have my uncommitted file
    And the branch "child_feature" has not been pushed to the repository
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE                                 | FILE NAME    |
      | main           | local and remote | main_commit                             | main_file    |
      | child-feature  | local            | feature_commit                          | feature_file |
      |                |                  | main_commit                             | main_file    |
      |                |                  | Merge branch 'main' into parent-feature |              |
      | parent-feature | local and remote | feature_commit                          | feature_file |
      |                |                  | main_commit                             | main_file    |
      |                |                  | Merge branch 'main' into parent-feature |              |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |


  Scenario: Providing '.' as the parent name
    When I run `git hack child-feature .`
    Then it runs the Git commands
      | BRANCH         | COMMAND                                      |
      | parent-feature | git fetch --prune                            |
      |                | git stash -u                                 |
      |                | git checkout main                            |
      | main           | git rebase origin/main                       |
      |                | git checkout parent-feature                  |
      | parent-feature | git merge --no-edit origin/parent-feature    |
      |                | git merge --no-edit main                     |
      |                | git push                                     |
      |                | git checkout -b child-feature parent-feature |
      | child-feature  | git stash pop                                |
    And I end up on the "child-feature" branch
    And I still have my uncommitted file
    And the branch "child_feature" has not been pushed to the repository
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE                                 | FILE NAME    |
      | main           | local and remote | main_commit                             | main_file    |
      | child-feature  | local            | feature_commit                          | feature_file |
      |                |                  | main_commit                             | main_file    |
      |                |                  | Merge branch 'main' into parent-feature |              |
      | parent-feature | local and remote | feature_commit                          | feature_file |
      |                |                  | main_commit                             | main_file    |
      |                |                  | Merge branch 'main' into parent-feature |              |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |
