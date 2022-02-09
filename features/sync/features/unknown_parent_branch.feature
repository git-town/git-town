@skipWindows
Feature: enter a parent branch name when prompted

  Background:
    Given the branches "alpha" and "beta"
    And the current branch is "beta"

  Scenario: choose the default branch name
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the parent branch of 'beta' | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | beta   | main   |

  Scenario: choose other branches
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                      | ANSWER        |
      | Please specify the parent branch of 'beta'  | [DOWN][ENTER] |
      | Please specify the parent branch of 'alpha' | [ENTER]       |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | alpha  | main   |
      | beta   | alpha  |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                     | ANSWER      |
      | Please specify the parent branch of 'beta' | [UP][ENTER] |
    Then the perennial branches are now "beta"

  Scenario: enter the parent for several branches
    When I run "git-town sync --all" and answer the prompts:
      | PROMPT                                      | ANSWER  |
      | Please specify the parent branch of 'alpha' | [ENTER] |
      | Please specify the parent branch of 'beta'  | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | alpha  | main   |
      | beta   | main   |
