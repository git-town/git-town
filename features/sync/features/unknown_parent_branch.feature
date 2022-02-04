@skipWindows
Feature: enter a parent branch name when prompted

  Background:
    Given my repo has the branches "one" and "two"
    And I am on the "two" branch

  Scenario: choose the default branch name
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                    | ANSWER  |
      | Please specify the parent branch of 'two' | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | two    | main   |

  Scenario: choose other branches
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                    | ANSWER        |
      | Please specify the parent branch of 'two' | [DOWN][ENTER] |
      | Please specify the parent branch of 'one' | [ENTER]       |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | one    | main   |
      | two    | one    |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                    | ANSWER      |
      | Please specify the parent branch of 'two' | [UP][ENTER] |
    Then the perennial branches are now "two"

  Scenario: enter the parent for several branches
    When I run "git-town sync --all" and answer the prompts:
      | PROMPT                                    | ANSWER  |
      | Please specify the parent branch of 'one' | [ENTER] |
      | Please specify the parent branch of 'two' | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | one    | main   |
      | two    | main   |
