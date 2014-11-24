Feature: git-pr when origin is unsupported

  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `git pr` while allowing errors


  Scenario:
    Then I get the error "Unsupported domain. Only Github and Bitbucket are currently supported"
