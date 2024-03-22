Feature: remove parent entries for perennial branches

  Scenario: parent branch entry for a perennial branch exists
    Given the local feature branches "feature-1" and "feature-2"
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = [ "feature-1" ]
      """
    When I run "git town config"
    Then this branch lineage exists now
      | BRANCH    | PARENT |
      | feature-2 | main   |
