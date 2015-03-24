Feature: git kill: killing the given feature branch (with open changes and without remote repo)

  (see ../with_tracking_branch/with_open_changes.feature)


  Background:
    Given I have feature branches named "feature" and "dead-feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH       | LOCATION | MESSAGE                              | FILE NAME        | FILE CONTENT   |
      | main         | local    | conflicting with uncommitted changes | conflicting_file | master content |
      | feature      | local    | good commit                          | good_file        |                |
      | dead-feature | local    | dead-end commit                      | unfortunate_file |                |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "conflicting_file" and content: "conflicting content"
    When I run `git kill dead-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                    |
      | feature | git branch -D dead-feature |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "conflicting_file" and content: "conflicting content"
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE                              | FILE NAME        |
      | main    | local    | conflicting with uncommitted changes | conflicting_file |
      | feature | local    | good commit                          | good_file        |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH  | COMMAND                                              |
      | feature | git branch dead-feature <%= sha 'dead-end commit' %> |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "conflicting_file" and content: "conflicting content"
    And the existing branches are
      | REPOSITORY | BRANCHES                    |
      | local      | main, dead-feature, feature |
    And I have the following commits
      | BRANCH       | LOCATION | MESSAGE                              | FILE NAME        |
      | main         | local    | conflicting with uncommitted changes | conflicting_file |
      | feature      | local    | good commit                          | good_file        |
      | dead-feature | local    | dead-end commit                      | unfortunate_file |
