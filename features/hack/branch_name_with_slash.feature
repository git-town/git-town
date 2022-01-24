Feature: git town-hack: does not error when branch name contains a forward slash

  Scenario: result
    Given I am on the "main" branch
    When I run "git-town hack my/feature"
    Then I am now on the "my/feature" branch
