@skipWindows
Feature: ask for missing parent

  Scenario: on feature branch without parent
    Given my repo has a branch "feature"
    And I am on the "feature" branch
    When I run "git-town diff-parent" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |
    And Git Town now knows this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
