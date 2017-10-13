Feature: git town-ship: errors when trying to ship the main branch

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given my repository has a feature branch named "feature"
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run `git-town ship main`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And my repository is still on the "feature" branch
    And my workspace still has my uncommitted file
