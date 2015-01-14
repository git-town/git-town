Feature: git-extract: errors if there are not extractable commits

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "The branch 'feature' has no extractable commits."
    And I am still on the "feature" branch
