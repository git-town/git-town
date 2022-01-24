Feature: git town-parent-diff: diffing the current feature branch

  As a user running parent-diff
  With no arguments
  On a feature branch that has no parent branch defined
  I should see a prompt to supply a parent branch
  So that the command can work as I expect

  @skipWindows
  Scenario: result
    Given my repo has a feature branch named "feature" with no parent
    And I am on the "feature" branch
    When I run "git-town diff-parent" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |
    And I am still on the "feature" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
