@skipWindows
Feature: ask for missing parent

  Scenario: feature branch without parent
    Given a branch "feature"
    And the current branch is "main"
    When I run "git-town diff-parent feature" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git diff main..feature |
    And this branch hierarchy exists now
      | BRANCH  | PARENT |
      | feature | main   |
