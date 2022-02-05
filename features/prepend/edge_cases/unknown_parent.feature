Feature: ask for missing parent information

  Scenario:
    Given my repo has a branch "old"
    And I am on the "old" branch
    When I run "git-town prepend new" and answer the prompts:
      | PROMPT                                    | ANSWER  |
      | Please specify the parent branch of 'new' | [ENTER] |
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git branch new main      |
      |        | git checkout new         |
    And Git Town now knows this branch hierarchy
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |
