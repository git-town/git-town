Feature: git rename-branch: renaming a non-feature branch without a remote repo (with open changes)

  (see ../with_remote_origin/with_tracking_branch/with_open_changes.feature)


  Background:
    Given my repo does not have a remote origin
    And I have a local branch named "production"
    And the following commits exist in my repository
      | BRANCH          | LOCATION | MESSAGE        |
      | main            | local    | main commit    |
      | production | local    | production commit |
    And I am on the "production" branch
    And I have an uncommitted file
    When I run `git rename-branch production renamed-production`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                                         |
      | production | git stash -u                                    |
      |                 | git checkout -b renamed-production production |
      | renamed-production | git branch -D production                   |
      |                 | git stash pop                                   |
    And I end up on the "renamed-production" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH          | LOCATION | MESSAGE        |
      | main            | local    | main commit    |
      | renamed-production | local    | production commit |
