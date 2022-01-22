Feature: Rendering Shell autocomplete definitions

  To use Git Town efficiently on the command line
  I want to enable shell autocompletion for Git Town commands.

  Scenario: verifying command output plausibility for fish autocompletion
    Given I run "git-town completions fish"
    Then it prints:
      """
      # fish completion for git-town
      """

  Scenario: verifying command output plausibility for Bash autocompletion
    Given I run "git-town completions bash"
    Then it prints:
      """
      # bash completion for git-town
      """

  Scenario: verifying command output plausibility for zsh autocompletion
    Given I run "git-town completions zsh"
    Then it prints:
      """
      # zsh completion for git-town
      """

  Scenario: verifying command output plausibility for PowerShell autocompletion
    Given I run "git-town completions powershell"
    Then it prints:
      """
      # powershell completion for git-town
      """
