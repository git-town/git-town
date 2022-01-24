Feature: git-town sync: syncing all branches syncs the tags

  As a developer using Git tags for release management
  I want my tags to be published whenever I sync all my branches
  So that I can do tagging work effectively on my local machine.

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
