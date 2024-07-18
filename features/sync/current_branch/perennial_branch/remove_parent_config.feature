Feature: remove parent entries for perennial branches

  Scenario: parent branch entry for a perennial branch exists
    Given a Git repo clone
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS |
      | feature-1 | feature | main   | local     |
      | feature-2 | feature | main   | local     |
    And the current branch is "feature-1"
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = [ "feature-2" ]
      """
    When I run "git town sync"
    Then it prints:
      """
      Removed parent entry for perennial branch "feature-2"
      """
    And this lineage exists now
      | BRANCH    | PARENT |
      | feature-1 | main   |
