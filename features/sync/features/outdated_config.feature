Feature: automatically remove outdated git-town configuration

  Scenario: automatically remove outdated branch hierarchy information
    Given I ran "git-town hack feature"
    And I ran "git checkout main"
    And I ran "git branch -D feature"
    When I run "git-town sync"
    Then no branch hierarchy exists now
