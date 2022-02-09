Feature: ask for missing parent branch information

  Scenario:
    Given my repo has a branch "feature"
    And I am on the "feature" branch
    When I run "git-town append new" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
