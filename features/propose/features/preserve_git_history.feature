Feature: preserve the previous Git branch

  Scenario:
    Given the feature branches "previous" and "current"
    And tool "open" is installed
    And the origin is "https://github.com/git-town/git-town.git"
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town propose"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
