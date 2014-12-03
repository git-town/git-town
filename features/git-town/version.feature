Feature: Version

  Scenario: Show version
    When I run `git town --version` while allowing errors
    Then the output should contain 'Git Town v0.4.1'