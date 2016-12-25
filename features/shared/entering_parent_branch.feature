Feature: Entering a parent branch name when prompted

  As a developer syncing a feature branch without information about its place in the branch hierarchy
  I want to be be able to enter the parent branch efficiently
  So that I am not slowed down much by the process of entering the parent branch.


  Background:
    Given I have feature branches named "feature-1" and "feature-2"
    And Git Town has no branch hierarchy information for "feature-1" and "feature-2"
    And I am on the "feature-2" branch


  Scenario: choosing the default branch name
    When I run `git town-sync` and press ENTER
    Then I see "Please specify the parent branch of feature-2"
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-2 | main   |


  Scenario: entering the number of the master branch
    When I run `git town-sync` and enter "1"
    Then I see "Please specify the parent branch of feature-2"
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-2 | main   |


  Scenario: entering the number of another branch
    When I run `git town-sync` and enter "2" and "1"
    Then I see "Please specify the parent branch of feature-2"
    And I see "Please specify the parent branch of feature-1"
    Then Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |


  Scenario: entering a wrong number
    When I run `git town-sync` and enter "5" and "1"
    Then I see "Please specify the parent branch of feature-2"
    And I see "invalid branch number"
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-2 | main   |


  Scenario: entering the name of the master branch
    When I run `git town-sync` and enter "main"
    Then I see "Please specify the parent branch of feature-2"
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-2 | main   |


  Scenario: entering the name of another branch
    When I run `git town-sync` and enter "feature-1" and "main"
    Then I see "Please specify the parent branch of feature-2"
    And I see "Please specify the parent branch of feature-1"
    Then Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |


  Scenario: entering a wrong name
    When I run `git town-sync` and enter "zonk" and "main"
    Then I see "Please specify the parent branch of feature-2"
    And I see "branch 'zonk' doesn't exist"
    Then Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-2 | main   |


  Scenario: entering self
    When I run `git town-sync` and enter "feature-2" and "main"
    Then I see "Please specify the parent branch of feature-2"
    And I see "'feature-2' cannot be the parent of itself"
    Then Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-2 | main   |


  Scenario: creating a loop
    When I run `git town-sync` and enter "feature-1" and "feature-2" and "main"
    Then I see "Please specify the parent branch of feature-2"
    And I see "Nested branch loop detected: 'feature-1' is an ancestor 'feature-2'"
    Then Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |
