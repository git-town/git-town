Feature: preserve the previous Git branch

  Scenario:
    Given the feature branches "previous" and "current"
    And my computer has the "open" tool installed
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town repo"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"
