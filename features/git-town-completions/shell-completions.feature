Feature: Rendering Shell autocomplete definitions

  As a shell user
  I want to use Git Town with autocompletion
  So that I can use this tool productively.

  Scenario: loading autocompletion in fish
    Given no fish autocompletion installed
    When I run "git-town completions fish | source"
    And when I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """

  Scenario: loading autocompletion in Bash
    Given no Bash autocompletion installed
    When I run "source <(git-town completions bash)"
    And when I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """

  Scenario: loading autocompletion in zsh
    Given no zsh autocompletion installed
    When I run "source <(git-town completions zsh)"
    And when I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """

  Scenario: loading autocompletion in PowerShell
    Given no fish autocompletion installed
    When I run "source <(git-town completions zsh)"
    And when I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """