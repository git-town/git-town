Feature: remove parent entries for perennial branches

  Scenario: parent branch entry for a perennial branch exists
    Given the current branch is a local feature branch "feature-1"
    Given the local feature branch "feature-2"
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = [ "feature-2" ]
      """
    When I run "git town park"
    Then this branch lineage exists now
      | BRANCH    | PARENT |
      | feature-1 | main   |
