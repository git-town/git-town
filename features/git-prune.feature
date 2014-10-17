Feature: Git Prune

  Scenario: on the main branch with a stale branch
    Given I have a feature branch named "feature" behind main
    And I am on the main branch
    When I run `git prune`
    Then I end up on the "main" branch
    Then the branch "feature" will be deleted


  Scenario: on the main branch with a remote stale branch
    Given my coworker has a feature branch named "feature" behind main
    And I am on the main branch
    When I run `git prune`
    Then I end up on the "main" branch
    Then the branch "feature" will be deleted


  Scenario: on a stale branch
    Given I have a feature branch named "feature" behind main
    And I am on the main branch
    When I run `git prune`
    Then I end up on the "main" branch
    Then the branch "feature" will be deleted


  Scenario: on the main branch with no stale branches
    Given I have a feature branch named "feature1" ahead of main
    And my coworker has a feature branch names "feature2" ahead of main
    And I am on the main branch
    When I run `git prune`
    Then I end up on the "main" branch
    And the branch "feature1" still exists
    And the branch "feature2" still exists
