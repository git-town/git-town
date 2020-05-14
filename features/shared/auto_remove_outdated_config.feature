Feature: automatically remove outdated git-town configuration

  As a developer using git-town
  I want my outdated configuration to be automatically removed
  So that my Git configuration isn’t littered with outdated entries.


  Scenario: automatically remove outdated branch hierarchy information
    Given I run "git-town hack feature"
    And I run "git checkout main"
    And I run "git branch -d feature"
    When I run "git-town sync"
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
