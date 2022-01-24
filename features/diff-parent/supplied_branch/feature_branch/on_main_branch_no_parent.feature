Feature: git town-parent-diff: diffing the current feature branch

  As a user running parent-diff
  With a supplied branch
  On the main branch
  I should see a prompt to identify a parent branch
  So that the command can work as I expect

  @skipWindows
  Scenario: result
    Given my repo has a feature branch named "feature" with no parent
    And I am on the "main" branch
    When I run "git-town diff-parent feature" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git diff main..feature |
    And I am still on the "main" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
