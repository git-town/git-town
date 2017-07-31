Feature: git town-ship: errors when trying to ship a perennial branch

  (see ../current_branch/on_perennial_branch.feature)


  Background:
    Given I have perennial branches named "qa" and "production"
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git-town ship production`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
    And I still have my uncommitted file
