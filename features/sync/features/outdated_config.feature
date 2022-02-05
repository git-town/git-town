Feature: automatically remove outdated git-town configuration

  Scenario: automatically remove outdated branch hierarchy information
    Given I run "git-town hack feature"
    And I run "git checkout main"
    And I run "git branch -d feature"
    When I run "git-town sync"
    And Git Town now knows about no branch hierarchy
