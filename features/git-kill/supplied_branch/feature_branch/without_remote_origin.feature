Feature: git kill: killing the given feature branch (without remote repo)

  (see ../with_tracking_branch/with_open_changes.feature)


  Background:
    Given my repo does not have a remote origin
    And I have local feature branches named "current-feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION | MESSAGE                | FILE NAME            | FILE CONTENT            |
      | main            | local    | main commit            | conflicting_file     | main content            |
      | current-feature | local    | current feature commit | current_feature_file | current feature content |
      | other-feature   | local    | other feature commit   | other_feature_file   | other feature content   |
    And I am on the "current-feature" branch
    And I have an uncommitted file with name: "conflicting_file" and content: "conflicting content"
    When I run `git kill other-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                     |
      | current-feature | git branch -D other-feature |
    And I am still on the "current-feature" branch
    And my workspace still has an uncommitted file with name: "conflicting_file" and content: "conflicting content"
    And the existing branches are
      | REPOSITORY | BRANCHES              |
      | local      | main, current-feature |
    And now I have the following commits
      | BRANCH          | LOCATION | MESSAGE                | FILE NAME            | FILE CONTENT            |
      | main            | local    | main commit            | conflicting_file     | main content            |
      | current-feature | local    | current feature commit | current_feature_file | current feature content |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the commands
      | BRANCH          | COMMAND                                                    |
      | current-feature | git branch other-feature <%= sha 'other feature commit' %> |
    And I am still on the "current-feature" branch
    And my workspace still has an uncommitted file with name: "conflicting_file" and content: "conflicting content"
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
    And now I have the following commits
      | BRANCH          | LOCATION | MESSAGE                | FILE NAME            | FILE CONTENT            |
      | main            | local    | main commit            | conflicting_file     | main content            |
      | current-feature | local    | current feature commit | current_feature_file | current feature content |
      | other-feature   | local    | other feature commit   | other_feature_file   | other feature content   |
