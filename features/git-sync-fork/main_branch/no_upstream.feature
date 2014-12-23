Feature: git-sync-fork without an upstream

  Scenario:
    When I run `git sync-fork` while allowing errors
    Then I get the error "Please add a remote 'upstream'"
