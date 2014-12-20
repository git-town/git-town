Feature: git kill: removing the current local feature branch (without open changes)

  As a developer on a local dead-end feature branch
  I want to be able to cleanly delete the current branch
  So that my workspace is ready to work on something else and my productivity remains high.


  Background:
    Given I have a feature branch named "good-feature"
    And I have a local feature branch named "dead-feature"
    And the following commits exist in my repository
      | branch       | location         | message         | file name        |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local            | dead-end commit | unfortunate_file |
    And I am on the local "dead-feature" branch
    When I run `git kill`


  Scenario: result
    Then I end up on the "main" branch
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |


  Scenario: Undoing a kill of a local feature branch
    When I run `git kill --undo`
    Then I end up on the "dead-feature" branch
    And the existing branches are
      | repository | branches                         |
      | local      | main, dead-feature, good-feature |
      | remote     | main, good-feature               |
    And I have the following commits
      | branch       | location         | message         | files            |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local            | dead-end commit | unfortunate_file |
