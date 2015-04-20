Feature: git rename-branch: renaming a feature branch with a tracking branch (with open changes)

  As a developer with a poorly named feature branch
  I want to be able to rename it safely in one easy step
  So that I can stay organized and remain productive


  Background:
    Given I have a feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE        |
      | main            | local and remote | main commit    |
      | current-feature | local and remote | feature commit |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `git rename-branch current-feature renamed-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                                         |
      | current-feature | git fetch --prune                               |
      |                 | git stash -u                                    |
      |                 | git checkout -b renamed-feature current-feature |
      | renamed-feature | git push -u origin renamed-feature              |
      |                 | git push origin :current-feature                |
      |                 | git branch -D current-feature                   |
      |                 | git stash pop                                   |
    And I end up on the "renamed-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE        |
      | main            | local and remote | main commit    |
      | renamed-feature | local and remote | feature commit |
