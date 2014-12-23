Feature: git-hack: errors when the branch exists locally (without open changes)

  As a developer trying to start a new feature on an already existing branch
  I should see an error telling me that the branch name is taken
  So that my feature branches are focussed, code reviews easy, and the team productivity remains high.


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
