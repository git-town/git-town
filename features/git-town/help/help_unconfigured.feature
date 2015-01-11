Feature: show help screen when Git Town is not configured

  (see ./help_configured.feature)


  Background:
    Given I haven't configured Git Town yet


  Scenario: git town with no flags
    When I run `git town`
    Then I see the git-town man page


  Scenario: git town, configured, with "help" subcommand
    When I run `git town help`
    Then I see the git-town man page
