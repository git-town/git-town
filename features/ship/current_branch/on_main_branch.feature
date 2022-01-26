Feature: git town-ship: errors when trying to ship the main branch

  Background:
    Given I am on the "main" branch
    When I run "git-town ship -m 'something done'"

  Scenario: result
    Then it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can be shipped
      """
    And I am still on the "main" branch
    And my repo now has the following commits
      | BRANCH | LOCATION |
