Feature: showing the current Git Town version

  Scenario: Using "--version" flag
    When I run `git town --version`
    Then I see "Git Town 0.4.1"

  Scenario: Using "version" flag
    When I run `git town version`
    Then I see "Git Town 0.4.1"
