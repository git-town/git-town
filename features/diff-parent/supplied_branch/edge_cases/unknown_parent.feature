@skipWindows
Feature: ask for missing parent

  Scenario: feature branch without parent
    Given a branch "feature"
    And I am on the "main" branch
    When I run "git-town diff-parent feature" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git diff main..feature |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
