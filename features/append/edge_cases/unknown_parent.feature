Feature: ask for missing parent branch information

  @debug @this
  Scenario:
    Given the current branch is "feature"
    When I run "git-town append new" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then this branch lineage exists now
      | BRANCH  | PARENT  |
      | feature | main    |
      | new     | feature |
