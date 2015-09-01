Feature: git-sync-fork without an upstream

  Scenario:
    When I run `git sync-fork`
    Then it runs no commands
    And I get the error "Please add a remote 'upstream'"
