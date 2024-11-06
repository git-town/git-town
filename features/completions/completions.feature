@smoke
Feature: shell autocompletion

  Scenario Outline:
    Given I am outside a Git repo
    When I run "git-town completions <SHELL>"
    Then Git Town prints:
      """
      # <SHELL> completion for git-town
      """

    Examples:
      | SHELL      |
      | fish       |
      | bash       |
      | zsh        |
      | powershell |
