Feature: cannot extract if a branch already exists with that name

  Background:
    Given I have feature branches named "feature" and "refactor"
    And the following commits exist in my repository
      | branch  | location | message            | file name        |
      | main    | remote   | remote main commit | remote_main_file |
      | feature | local    | feature commit     | feature_file     |
      | feature | local    | refactor commit    | refactor_file    |
    And I am on the "feature" branch
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then I get the error "A branch named 'refactor' already exists"
    And I am still on the "feature" branch
