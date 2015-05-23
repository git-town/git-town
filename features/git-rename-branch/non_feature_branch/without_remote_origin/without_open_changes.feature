Feature: git rename-branch: renaming a non-feature branch without a remote repo (without open changes)

  (see ../with_remote_origin/with_tracking_branch/with_open_changes.feature)


  Background:
    Given my repo does not have a remote origin
    And I have local branches named "qa" and "production"
    And my non-feature branches are configured as "qa" and "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION | MESSAGE           |
      | main       | local    | main commit       |
      | production | local    | production commit |
      | qa         | local    | qa commit         |
    And I am on the "production" branch


  Scenario: error when trying to rename
    When I run `git rename-branch production renamed-production`
    Then it runs no Git commands
    And I get the error "The branch 'production' is not a feature branch."
    And I get the error "Run 'git rename-branch production renamed-production -f' to force the rename, then reconfigure git-town on any other clones of this repo."


  Scenario: result
    When I run `git rename-branch production renamed-production -f`
    Then it runs the Git commands
      | BRANCH             | COMMAND                                       |
      | production         | git checkout -b renamed-production production |
      | renamed-production | git branch -D production                      |
    And I end up on the "renamed-production" branch
    And my non-feature branches are now configured as "qa" and "renamed-production"
    And I have the following commits
      | BRANCH             | LOCATION | MESSAGE           |
      | main               | local    | main commit       |
      | qa                 | local    | qa commit         |
      | renamed-production | local    | production commit |
