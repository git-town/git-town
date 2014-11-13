Feature: Git Kill

  Background:
    Given I have a feature branch named "good-feature"


  Scenario: Killing a feature branch without open changes
    Given I am on the "stupid-feature" branch
    When I run `git kill`
    Then I end up on the "main" branch
    And the branch "stupid-feature" is deleted everywhere
    And the branch "good-feature" still exists


  Scenario: Killing a feature branch with open changes
    Given I am on the "stupid-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`
    Then I end up on the "main" branch
    And the branch "stupid-feature" is deleted everywhere
    And the branch "good-feature" still exists
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: Does not kill the main branch
    Given I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill` while allowing errors
    Then I get the error "You can only kill feature branches"
    And I am still on the "main" branch
    And the branch "good-feature" still exists
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: Does not kill a non-feature branch
    Given non-feature branch configuration "qa"
    And I am on the "qa" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill` while allowing errors
    Then I get the error "You can only kill feature branches"
    And I am still on the "qa" branch
    And the branch "good-feature" still exists
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: Undoing a kill without open changes
    Given I am on the "unfortunate-feature" branch
    When I run `git kill`
    And I run `git kill --undo`
    Then I end up on the "unfortunate-feature" branch
    And the existing branches are
      | repository | branches                                |
      | local      | main, unfortunate-feature, good-feature |
      | remote     | main, unfortunate-feature, good-feature |


  Scenario: Undoing a kill with open changes
    Given I am on the "unfortunate-feature" branch
    And the following commits exist in my repository
      | branch              | location         | message            | file name        |
      | good-feature        | local and remote | good commit        | good_file        |
      | unfortunate-feature | local and remote | unfortunate commit | unfortunate_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`
    And I run `git kill --undo`
    Then I end up on the "unfortunate-feature" branch
    And the existing branches are
      | repository | branches                                |
      | local      | main, unfortunate-feature, good-feature |
      | remote     | main, unfortunate-feature, good-feature |
    And I have the following commits
      | branch              | location         | message            | files            |
      | good-feature        | local and remote | good commit        | good_file        |
      | unfortunate-feature | local and remote | unfortunate commit | unfortunate_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: Undoing a kill of a local feature branch
    Given I am on the local "unfortunate-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`
    And I run `git kill --undo`
    Then I end up on the "unfortunate-feature" branch
    And the existing branches are
      | repository | branches                                |
      | local      | main, unfortunate-feature, good-feature |
      | remote     | main, good-feature                      |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
