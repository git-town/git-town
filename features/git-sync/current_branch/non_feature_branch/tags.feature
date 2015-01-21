Feature: git sync: syncing the current non-feature branch pushes tags to the remote

  As a developer syncing a non-feature branch
  I want my tags to be published
  So that I can do tagging work effectively on my local machine and have more time for other work.


  Background:
    Given I have branches named "qa" and "production"
    And my non-feature branches are configured as "qa" and "production"
    And I am on the "production" branch
    And I add a local tag "v1.0"
    When I run `git sync`


  Scenario: result
    Then tag "v1.0" has been pushed to the remote
