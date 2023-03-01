@skipWindows
Feature: update the parent of a feature branch

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"

  Scenario: select the default branch (current parent)
    When I run "git-town set-parent" and answer the prompts:
      | PROMPT                                      | ANSWER  |
      | Please specify the parent branch of 'child' | [ENTER] |
    And the initial branch hierarchy exists

  Scenario: select another branch
    When I run "git-town set-parent" and answer the prompts:
      | PROMPT                                      | ANSWER      |
      | Please specify the parent branch of 'child' | [UP][ENTER] |
    Then this branch hierarchy exists now
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town set-parent" and answer the prompts:
      | PROMPT                                      | ANSWER          |
      | Please specify the parent branch of 'child' | [UP][UP][ENTER] |
    Then the perennial branches are now "child"
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | parent | main   |
