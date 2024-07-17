@skipWindows
Feature: Prepopulate title and body

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given the current branch is "feature"
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    And tool "open" is installed

  Scenario: provide title and body via CLI
    When I run "git-town propose --title=my_title --body=my_body"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title&body=my_body
      """

  Scenario: provide title via CLI
    When I run "git-town propose --title=my_title"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title
      """

  Scenario: provide body via CLI
    When I run "git-town propose --body=my_body"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&body=my_body
      """

  Scenario: provide title via CLI and body via file
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

  Scenario: non-existing body file
    When I run "git-town propose --body-file zonk.txt"
    Then it prints the error:
      """
      Error: open zonk.txt: no such file or directory
      """

  Scenario: provide title via CLI and body via STDIN
    When I pipe the following text into "git-town propose --body-file -":
      """
      Proposal
      body
      text
      """
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1&body=Proposal%0Abody%0Atext
      """
