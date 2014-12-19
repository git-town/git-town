Feature: git-hack errors when branch exists on local (without open changes)

  Background:
    Given I have a feature branch named "existing_feature"
    And I am on the main branch
    When I run `git hack existing_feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
