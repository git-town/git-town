Feature: display the parent of a perennial branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | perennial | perennial |        | local, origin |
    And the current branch is "perennial"
    When I run "git-town config get-parent"

  Scenario: result
    Then it runs no commands
    And it prints no output
