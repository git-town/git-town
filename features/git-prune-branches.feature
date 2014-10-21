Feature: Git Prune Branches

  Scenario: on the main branch with a feature branch behind main
    Given I have a feature branch named "old_feature" behind main
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And the branch "old_feature" is deleted


  Scenario: on the main branch with a remote feature branch behind main
    Given my coworker has a feature branch named "old_feature" behind main
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And the branch "old_feature" is deleted


  Scenario: on a feature branch ahead of main with a feature branch behind main
    Given I have a feature branch named "old_feature" behind main
    And I have a feature branch named "new_feature" ahead of main
    And I am on the "new_feature" branch
    When I run `git prune-branches`
    Then I end up on the "new_feature" branch
    And the branch "old_feature" is deleted


  Scenario: on a feature branch behind main
    Given I have a feature branch named "old_feature" behind main
    And I am on the "old_feature" branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And the branch "old_feature" is deleted


  Scenario: on the main branch with feature branches ahead of main
    Given I have a feature branch named "feature1" ahead of main
    And my coworker has a feature branch named "feature2" ahead of main
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And the branch "feature1" still exists
    And the branch "feature2" still exists


  Scenario: on the main branch with a non-feature branch behind main
    Given I have a non-feature branch "production" behind main
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And the branch "production" still exists
