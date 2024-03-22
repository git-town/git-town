Feature: remove parent entries for perennial branches

  Scenario: parent branch entry for a perennial branch exists
    Given the local feature branches "feature-1" and "feature-2"
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = [ "feature-1" ]
      """
    When I run "git town append feature-3"
    Then it prints:
      """
      Removed parent entry for perennial branch "feature-1"
      """
    And this branch lineage exists now
      | BRANCH    | PARENT |
      | feature-2 | main   |
      | feature-3 | main   |
