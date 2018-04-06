Feature: automatically remove outdated git-town configuration

  As a developer using git-town
  I want my outdated configuration to be automatically removed
  So that my Git configuration isn’t littered with outdated entries.


  Scenario: automatically remove outdated branch hierarchy information
    Given I run `git-town hack feature; git checkout main; git branch -d feature`
    When I run `git-town sync`
    Then Git Town has no branch hierarchy information
