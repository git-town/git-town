@skipWindows
Feature: strip colors

  Scenario: colors are stripped from the output of git commands run internally
    Given Git Town is not configured
    And Git setting "color.ui" is "always"
    And the current branch is "main"
    When I run "git-town hack new" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then this branch hierarchy exists now
      | BRANCH | PARENT |
      | new    | main   |
