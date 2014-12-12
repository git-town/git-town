Feature: git-repo when origin is on Bitbucket over HTTPS

  Background:
    Given my remote origin is on Bitbucket through HTTPS
    When I run `git repo`


  Scenario:
    Then I see a browser window for my repository homepage on Bitbucket
