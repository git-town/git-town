Feature: remove a branch from the perennial branches configuration


  Background:
    Given my repo has the branches "staging" and "qa"
    And the perennial branches are configured as "staging" and "qa"


  @skipWindows
  Scenario: removing a branch that is a perennial branch
    When I run "git-town perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now configured as "qa"
