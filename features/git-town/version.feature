Feature: Version

  Scenario: Using "--version" flag
    When I run `git town --version` while allowing errors
    Then the output should contain "Git Town 0.4.1"

  Scenario: Using "version" flag
    When I run `git town version` while allowing errors
    Then the output should contain "Git Town 0.4.1"
