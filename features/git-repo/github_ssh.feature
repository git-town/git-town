Feature: git-repo when origin is on GitHub over SSH

  Background:
    Given my remote origin is on GitHub through SSH
    When I run `git repo`


  Scenario: result
    Then I see a browser window for my repository homepage on GitHub
