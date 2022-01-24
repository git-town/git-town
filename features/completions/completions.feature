Feature: shell autocompletion

  Scenario: fish autocompletion
    Given I run "git-town completions fish"
    Then it prints:
      """
      # fish completion for git-town
      """

  Scenario: bash autocompletion
    Given I run "git-town completions bash"
    Then it prints:
      """
      # bash completion for git-town
      """

  Scenario: zsh autocompletion
    Given I run "git-town completions zsh"
    Then it prints:
      """
      # zsh completion for git-town
      """

  Scenario: powershell autocompletion
    Given I run "git-town completions powershell"
    Then it prints:
      """
      # powershell completion for git-town
      """
