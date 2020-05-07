Feature: git town-parent-diff: diffing the current feature branch

    As a user running parent-diff
    With a supplied branch that matches my current branch
    On a branch that has no parent branch defined
    I should see a prompt to supply a parent branch
    So that the command can work as I expect


  Background:
    Given my repository has a feature branch named "feature" with no parent
    And the following commits exist in my repository
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
    And I am on the "feature" branch
    And my workspace has an uncommitted file


  Scenario:
    When I run "git-town diff-parent feature" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
