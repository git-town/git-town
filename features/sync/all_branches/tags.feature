Feature: git-town sync: syncing all branches syncs the tags

  Background:
    Given my repo has the following tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | remote-tag | remote   |
    And I am on the "main" branch
    When I run "git-town sync --all"

  Scenario: result
    Then my repo now has the following tags
      | NAME       | LOCATION      |
      | local-tag  | local, remote |
      | remote-tag | local, remote |
