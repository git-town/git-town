Feature: git town-ship: errors when trying to ship a perennial branch

  (see ../current_branch/on_perennial_branch.feature)


  Background:
    Given my repository has the perennial branches "qa" and "production"
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town ship production`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And it prints the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
