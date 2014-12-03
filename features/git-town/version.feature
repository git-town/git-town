Feature: Version

  Scenario: Show version
    Given I am using Git Town version "3.141592"
    When I run `git town --version` while allowing errors
    Then the output should contain 'Git Town v<#GIT_TOWN_VERSION=3.141592#>'