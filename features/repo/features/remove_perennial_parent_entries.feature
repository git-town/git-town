Feature: remove parent entries for perennial branches

  Scenario: parent branch entry for a perennial branch exists
    Given tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And the current branch is a local feature branch "feature-1"
    And the local feature branch "feature-2"
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = [ "feature-2" ]
      """
    When I run "git town repo"
    Then this branch lineage exists now
      | BRANCH    | PARENT |
      | feature-1 | main   |
