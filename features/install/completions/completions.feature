Feature: shell autocompletion

  Scenario Outline:
    Given I run "git-town install completions <SHELL>"
    Then it prints:
      """
      # <SHELL> completion for git-town
      """

    Examples:
      | SHELL      |
      | fish       |
      | bash       |
      | zsh        |
      | powershell |
