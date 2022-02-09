@skipWindows
Feature: ask for missing parent

  Scenario: on feature branch without parent
    Given the current branch is "feature"
    When I run "git-town diff-parent" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |
    And this branch hierarchy exists now
      | BRANCH  | PARENT |
      | feature | main   |
