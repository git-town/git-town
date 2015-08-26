Feature: git hack: does not error when branch name contains a slash

  As a developer trying to create a branch name that contains a slash
  I don't want to see a warning about an invalid key
  So that the tool doesn't get in my way


  Scenario: result
    Given I am on the "main" branch
    When I run `git hack my/feature`
    Then I end up on the "my/feature" branch
