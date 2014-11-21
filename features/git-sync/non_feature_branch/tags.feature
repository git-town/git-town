Feature: Git Sync: syncing a non-feature branch pushes tags to the remote


  Scenario: Tags
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    And I add a local tag "v1.0"
    When I run `git sync`
    Then tag "v1.0" has been pushed to the remote
