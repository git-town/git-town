Feature: ask for missing parent information

  Scenario:
    Given my repo has a branch "feature"
    And I am on the "feature" branch
    When I run "git-town prepend new" and answer the prompts:
      | PROMPT                                    | ANSWER  |
      | Please specify the parent branch of 'new' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
      |         | git branch new main      |
      |         | git checkout new         |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | new    |
      | new     | main   |
