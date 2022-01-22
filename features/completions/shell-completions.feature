Feature: Rendering Shell autocomplete definitions

  To use Git Town efficiently
  When on the command line
  I want autocompletion for Git Town commands.

  Scenario: fish autocompletion
    Given I run "git-town completions fish"
    Then it prints:
      """
      # fish completion for git-town
      """

  Scenario: Bash autocompletion
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

  Scenario: PowerShell autocompletion
    Given I run "git-town completions powershell"
    Then it prints:
      """
      # powershell completion for git-town
      """
