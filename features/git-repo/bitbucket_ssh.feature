Feature: git-repo when origin is on Bitbucket over SSH

  Background:
    Given my remote origin is on Bitbucket through SSH
    When I run `git repo`


  Scenario: result
    Then I see a browser window for my repository homepage on Bitbucket
