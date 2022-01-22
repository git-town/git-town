Feature: Viewing changes made by the current feature branch

  To know whether my local branch is correctly set up
  When working with nested feature branches
  I want to see the changes my current branch makes.

  Scenario: known parent branch
    Given my repo has a feature branch named "feature-1"
    And my repo has a feature branch named "feature-2" as a child of "feature-1"
    And I am on the "feature-2" branch
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | feature-2 | git diff feature-1..feature-2 |

  Scenario: unknown parent branch
    Given my repo has a feature branch named "feature" with no parent
    And I am on the "feature" branch
    When I run "git-town diff-parent" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
