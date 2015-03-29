Feature: git kill: killing the given feature branch when on it (with open changes and without remote repo)

  (see ../with_tracking_branch/with_open_changes.feature)


  Background:
    Given I have feature branches named "current-feature" and "other-feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH          | LOCATION | MESSAGE                | FILE NAME            | FILE CONTENT            |
      | current-feature | local    | current feature commit | current_feature_file | current feature content |
      | other-feature   | local    | other feature commit   | other_feature_file   | other feature content   |
    And I am on the "current-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill current-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                                |
      | current-feature | git add -A                             |
      |                 | git commit -m 'WIP on current-feature' |
      |                 | git checkout main                      |
      | main            | git branch -D current-feature          |
    And I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
    And now I have the following commits
      | BRANCH        | LOCATION | MESSAGE              | FILE NAME          |
      | other-feature | local    | other feature commit | other_feature_file |
    And now I have the following committed files
      | BRANCH        | NAME               | CONTENT               |
      | other-feature | other_feature_file | other feature content |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH          | COMMAND                                                        |
      | main            | git branch current-feature <%= sha 'WIP on current-feature' %> |
      |                 | git checkout current-feature                                   |
      | current-feature | git reset <%= sha 'current feature commit' %>                  |
    And I end up on the "current-feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
    And now I have the following commits
      | BRANCH          | LOCATION | MESSAGE                | FILE NAME            |
      | current-feature | local    | current feature commit | current_feature_file |
      | other-feature   | local    | other feature commit   | other_feature_file   |
    And now I have the following committed files
      | BRANCH          | NAME                 | CONTENT                 |
      | current-feature | current_feature_file | current feature content |
      | other-feature   | other_feature_file   | other feature content   |
