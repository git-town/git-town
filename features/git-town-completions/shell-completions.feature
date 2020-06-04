Feature: Rendering Shell autocomplete definitions

  As a shell user
  I want to use Git Town with autocompletion
  So that I can use this tool productively.

  Scenario: loading autocompletion in fish
    Given a fish shell
    And I run "git-town completions fish | source"
    And then I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """

  Scenario: loading autocompletion in Bash
    Given a Bash shell
    And I run "source <(git-town completions bash)"
    And then I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """

  Scenario: loading autocompletion in zsh
    Given a zsh shell
    And I run "source <(git-town completions zsh)"
    And then I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """

  Scenario: loading autocompletion in PowerShell
    Given a PowerShell
    And I run "source <(git-town completions powershell)"
    And then I type "git-town <TAB>"
    Then it prints:
      """
      ????
      """