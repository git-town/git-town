Feature: shell autocompletion

  Scenario Outline:
    When I run "git-town completions <SHELL>"
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
