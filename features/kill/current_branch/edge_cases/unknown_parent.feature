Feature: ask for missing parent branch information

  Scenario:
    Given the current branch is "feature"
    When I run "git-town kill feature" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git branch -d feature    |
    And no lineage exists now
