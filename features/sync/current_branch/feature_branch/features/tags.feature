Feature: syncing a feature branch pulls tags

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
