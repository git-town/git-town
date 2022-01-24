Feature: git-town sync: syncing a feature branch pulls tags

  As a developer using Git tags for release management
  I want that tags are pulled automatically for me whenever I sync
  So that my local workspace has the same tags that exist on the remote

  Background:
    Given my repo has a feature branch named "feature"
    And my repo has the following tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | remote-tag | remote   |
    And I am on the "feature" branch
    And I run "git-town sync"

  Scenario: result
    Then my repo now has the following tags
      | NAME       | LOCATION      |
      | local-tag  | local         |
      | remote-tag | local, remote |
