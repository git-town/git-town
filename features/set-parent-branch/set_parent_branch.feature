@skipWindows
Feature: update the parent of a feature branch

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"

  Scenario: select the default branch (current parent)
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                      | ANSWER  |
      | Please specify the parent branch of 'child' | [ENTER] |
    And Git Town is still aware of the initial branch hierarchy

  Scenario: select another branch
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                      | ANSWER      |
      | Please specify the parent branch of 'child' | [UP][ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | child  | main   |
      | parent | main   |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                      | ANSWER          |
      | Please specify the parent branch of 'child' | [UP][UP][ENTER] |
    Then the perennial branches are now "child"
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | parent | main   |
