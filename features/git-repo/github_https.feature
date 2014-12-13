Feature: git-repo when origin is on GitHub over HTTPS

  Background:
    Given my remote origin is on GitHub through HTTPS
    When I run `git repo`


  Scenario: result
    Then I see a browser window for my repository homepage on GitHub
