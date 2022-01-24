Feature: git-town sync: syncing the current perennial branch syncs the tags

  As a developer using Git tags for release management
  I want my tags to be published whenever I sync a perennial branch
  So that I can do tagging work effectively on my local machine.

  Background:
    Given my repo has the perennial branches "production" and "qa"
    And I am on the "production" branch
    And my repo has the following tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | remote-tag | remote   |
    When I run "git-town sync"

  Scenario: result
    Then my repo now has the following tags
      | NAME       | LOCATION      |
      | local-tag  | local, remote |
      | remote-tag | local, remote |
