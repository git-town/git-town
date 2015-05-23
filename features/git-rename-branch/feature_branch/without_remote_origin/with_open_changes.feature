Feature: git rename-branch: renaming a feature branch without a remote repo (with open changes)

  (see ../with_remote_origin/with_tracking_branch/with_open_changes.feature)


  Background:
    Given my repo does not have a remote origin
    And I have a local feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION | MESSAGE        |
      | main            | local    | main commit    |
      | current-feature | local    | feature commit |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `git rename-branch current-feature renamed-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                                         |
      | current-feature | git stash -u                                    |
      |                 | git checkout -b renamed-feature current-feature |
      | renamed-feature | git branch -D current-feature                   |
      |                 | git stash pop                                   |
    And I end up on the "renamed-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH          | LOCATION | MESSAGE        |
      | main            | local    | main commit    |
      | renamed-feature | local    | feature commit |
