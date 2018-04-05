Feature: outdated git-town configuration is automatically removed

  As a developer using git town
  I want to git town to automatically clean up outdated config
  So that my config file doesn't endlessly grow


  Scenario:
    Given I run `git-town hack feature; git checkout main; git branch -d feature`
    When I run `git-town sync`
    Then Git Town has no branch hierarchy information
