@skipWindows
Feature: strip colors

  Scenario: colors are stripped from the output of git commands run internally
    Given Git Town is not configured
    And the "color.ui" setting is "always"
    And I am on the "main" branch
    When I run "git-town hack new-feature" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH      | PARENT |
      | new-feature | main   |
