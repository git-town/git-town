@skipWindows
Feature: update the parent of a nested feature branch

  Background:
    Given my repo has a feature branch "parent-feature"
    And my repo has a feature branch "child-feature" as a child of "parent-feature"
    And I am on the "child-feature" branch

  Scenario: select the default branch (current parent)
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                              | ANSWER  |
      | Please specify the parent branch of 'child-feature' | [ENTER] |
    And Git Town still has the original branch hierarchy

  Scenario: select another branch
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                              | ANSWER      |
      | Please specify the parent branch of 'child-feature' | [UP][ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | child-feature  | main   |
      | parent-feature | main   |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                              | ANSWER          |
      | Please specify the parent branch of 'child-feature' | [UP][UP][ENTER] |
    Then the perennial branches are now "child-feature"
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | parent-feature | main   |
