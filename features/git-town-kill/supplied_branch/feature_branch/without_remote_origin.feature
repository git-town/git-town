Feature: git town-kill: killing the given feature branch (without remote repo)

  (see ../with_tracking_branch/with_open_changes.feature)


  Background:
    Given my repo does not have a remote origin
    And my repository has the local feature branches "current-feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION | MESSAGE                | FILE NAME            |
      | main            | local    | main commit            | conflicting_file     |
      | current-feature | local    | current feature commit | current_feature_file |
      | other-feature   | local    | other feature commit   | other_feature_file   |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file with name: "conflicting_file" and content: "conflicting content"
    When I run `git-town kill other-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                     |
      | current-feature | git add -A                  |
      |                 | git stash                   |
      |                 | git branch -D other-feature |
      |                 | git stash pop               |
    And I am still on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES              |
      | local      | main, current-feature |
    And now my repository has the following commits
      | BRANCH          | LOCATION | MESSAGE                | FILE NAME            |
      | main            | local    | main commit            | conflicting_file     |
      | current-feature | local    | current feature commit | current_feature_file |


  Scenario: undoing the kill
    When I run `git-town undo`
    Then it runs the commands
      | BRANCH          | COMMAND                                                    |
      | current-feature | git add -A                                                 |
      |                 | git stash                                                  |
      |                 | git branch other-feature <%= sha 'other feature commit' %> |
      |                 | git stash pop                                              |
    And I am still on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
    And my repository is left with my original commits
