Feature: print URL if no tool to open browsers is available

  As a developer using Git Town on a remote machine via SSH
  I want that the tool prints the URL it would have opened
  So that I can copy-and-paste the URL to open it on my local computer.

  Background:
    Given my repository has a feature branch named "feature"
    And my repo's origin is "git@github.com:git-town/git-town"
    And my computer has no tool to open browsers installed
    And I am on the "feature" branch
    When I run "git-town new-pull-request"

  Scenario: result
    Then it prints
      """
      Cannot find a browser to open https://github.com/git-town/git-town/compare/feature?expand=1
      """

  Scenario: undo
    When I run "git-town undo"
    Then it prints
      """
      foo
      """
