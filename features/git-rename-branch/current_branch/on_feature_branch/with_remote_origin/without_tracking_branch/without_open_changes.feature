Feature: git rename-branch: renaming a feature branch without a tracking branch (without open changes)

  (see ../with_tracking_branch/with_open_changes.feature)

  Background:
    Given I have a local feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE        |
      | main            | local and remote | main commit    |
      | current-feature | local            | feature commit |
    And I am on the "current-feature" branch
    When I run `git rename-branch current-feature renamed-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                                         |
      | current-feature | git fetch --prune                               |
      |                 | git checkout -b renamed-feature current-feature |
      | renamed-feature | git branch -D current-feature                   |
    And I end up on the "renamed-feature" branch
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE        |
      | main            | local and remote | main commit    |
      | renamed-feature | local            | feature commit |
