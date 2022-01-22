Feature: git town parent-diff on a feature branch

  To know whether my branch setup is correct
  When working with nested feature branches
  I want to see the changes a feature branch makes.

  Scenario: known parent branch
    Given my repo has a feature branch named "feature-1"
    And my repo has a feature branch named "feature-2" as a child of "feature-1"
    And I am on the "feature-2" branch
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | feature-2 | git diff feature-1..feature-2 |
    And I am still on the "feature-2" branch

  Scenario: unknown parent branch
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
