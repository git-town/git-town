Feature: Git Prune Branches

  Scenario: on the main branch with a stale branch
    Given I have a feature branch named "old_feature"
    And the following commit exists in my repository
      | branch  | message     | file name |
      | main    | main commit | main_file |
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    Then all branches are now synchronized
    Then the branch "old_feature" will be deleted


  Scenario: on a feature branch with a stale branch
    Given I have a feature branch named "old_feature"
    And I have a feature branch named "new_feature"
    And the following commit exists in my repository
      | branch      | message            | file name        |
      | main        | main commit        | main_file        |
      | new_feature | new feature commit | new_feature_file |
    And I am on the "new_feature" branch
    When I run `git prune-branches`
    Then I end up on the "new_feature" branch
    And the branch "old_feature" will be deleted


  Scenario: on the main branch with a remote stale branch
    Given my coworker has a feature branch named "old_feature"
    And the following commit exists in my repository
      | branch  | message     | file name |
      | main    | main commit | main_file |
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    Then the branch "old_feature" will be deleted


  Scenario: on a stale branch
    Given I have a feature branch named "old_feature"
    And the following commit exists in my repository
      | branch  | message     | file name |
      | main    | main commit | main_file |
    And I am on the "old_feature" branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    Then the branch "old_feature" will be deleted


  Scenario: on the main branch with no stale branches
    Given I have a feature branch named "feature1"
    And my coworker has a feature branch named "feature2"
    And the following commits exist in my repository
      | branch   | message   | file name     |
      | feature1 | my commit | feature1 file |
    And the following commits exist in Charlie's repository
      | branch   | message   | file name     |
      | feature2 | my commit | feature2 file |
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And the branch "feature1" still exists
    And the branch "feature2" still exists


  Scenario: on the main branch with non feature branch
    Given I have a branch named "production"
    And non-feature branch configuration "qa, production"
    And the following commit exists in my repository
      | branch  | message     | file name |
      | main    | main commit | main_file |
    And I am on the main branch
    When I run `git prune-branches`
    Then I end up on the "main" branch
    And the branch "production" still exists
