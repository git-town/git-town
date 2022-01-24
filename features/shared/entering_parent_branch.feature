@skipWindows
Feature: Entering a parent branch name when prompted

  As a developer syncing a feature branch without information about its place in the branch hierarchy
  I want to be be able to enter the parent branch efficiently
  So that I am not slowed down much by the process of entering the parent branch.

  Background:
    Given my repo has the branches "feature-1" and "feature-2"
    And I am on the "feature-2" branch

  Scenario: choosing the default branch name
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                          | ANSWER  |
      | Please specify the parent branch of 'feature-2' | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-2 | main   |

  Scenario: choosing other branches
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                          | ANSWER        |
      | Please specify the parent branch of 'feature-2' | [DOWN][ENTER] |
      | Please specify the parent branch of 'feature-1' | [ENTER]       |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |

  Scenario: choosing "<none> (make a perennial branch)"
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                          | ANSWER      |
      | Please specify the parent branch of 'feature-2' | [UP][ENTER] |
    Then the perennial branches are now configured as "feature-2"
