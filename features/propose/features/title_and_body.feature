@skipWindows
Feature: Prepopulate title and body

  Background:
    Given tool "open" is installed

  Scenario: provide title and body via CLI
    Given the current branch is a feature branch "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    When I run "git-town propose --title=my_title --body=my_body"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title&body=my_body
      """

  Scenario: provide title via CLI
    Given the current branch is a feature branch "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    When I run "git-town propose --title=my_title"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title
      """

  Scenario: provide body via CLI
    Given the current branch is a feature branch "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    When I run "git-town propose --body=my_body"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&body=my_body
      """

  Scenario: provide title via CLI and body via file
    Given the current branch is a feature branch "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    And file "body.txt" with content
      """
      Proposal
      body
      text!
      """
    When I run "git-town propose --body-file=body.txt"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&body=Proposal%0Abody%0Atext%21
      """
