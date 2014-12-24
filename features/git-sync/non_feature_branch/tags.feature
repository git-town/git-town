Feature: git sync: syncing a non-feature branch pushes tags to the remote

  As a developer syncing a non-feature branch
  I want my tags to be published
  So that I can do tagging work effectively on my local machine and have more time for other work.

  Scenario: Tags
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    And I add a local tag "v1.0"
    When I run `git sync`
    Then tag "v1.0" has been pushed to the remote
