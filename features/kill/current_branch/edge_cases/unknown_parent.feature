Feature: ask for missing parent branch information

  Scenario:
    Given a branch "feature"
    And I am on the "feature" branch
    When I run "git-town kill feature" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git branch -D feature    |
    And Git Town is now aware of no branch hierarchy
