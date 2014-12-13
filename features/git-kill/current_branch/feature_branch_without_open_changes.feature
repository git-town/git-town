Feature: git kill: removing the current feature branch (without open changes)

  As a developer being on a dead-end feature branch
  I want to be able to cleanly delete the current branch everywhere
  So that my workspace is ready to work on something else and my productivity remains high.


  Background:
    Given I have feature branches named "good-feature" and "unfortunate"
    And the following commits exist in my repository
      | branch       | location         | message            | file name        |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local and remote | unfortunate commit | unfortunate_file |
    And I am on the "unfortunate" branch
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


  Scenario: Undoing the kill
    When I run `git kill --undo`
    Then I end up on the "unfortunate" branch
    And the existing branches are
      | repository | branches                        |
      | local      | main, unfortunate, good-feature |
      | remote     | main, unfortunate, good-feature |
    And I have the following commits
      | branch       | location         | message            | files            |
      | good-feature | local and remote | good commit        | good_file        |
      | unfortunate  | local and remote | unfortunate commit | unfortunate_file |
