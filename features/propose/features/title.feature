@skipWindows
Feature: Prepopulate title and body

  Background:
    Given tool "open" is installed

  Scenario Outline: provide title and body
    Given the current branch is a feature branch "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    When I run "git-town propose --title=my_title --body=my_body"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title&body=my_body
      """

  Scenario Outline: provide title only body
    Given the current branch is a feature branch "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    When I run "git-town propose --title=my_title"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title
      """

  Scenario Outline: provide body only
    Given the current branch is a feature branch "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    When I run "git-town propose --body=my_body"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&body=my_body
      """
