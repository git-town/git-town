Feature: preserve the previous Git branch

  Scenario:
    Given the feature branches "previous" and "current"
    And my computer has the "open" tool installed
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town repo"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
