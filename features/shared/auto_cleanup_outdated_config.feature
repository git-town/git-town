Feature: outdated git-town configuration is automatically removed

  Scenario:
    Given I run `git-town hack feature; git checkout main; git branch -d feature`
    When I run `git-town sync --debug`
    Then Git Town has no branch hierarchy information
