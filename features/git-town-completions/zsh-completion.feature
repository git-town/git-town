Feature: Installing Zsh Shell autocomplete definitions

  As a Zsh shell user
  I want to use Git Town with autocompletion
  So that I can use this tool productively.


  Scenario: without existing zsh autocompletion folder
    Given my computer has no zsh autocompletion file
    When I run "git-town completions zsh"
    Then it prints:
      """
      Zsh autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with empty zsh autocompletion folder
    Given my computer has an empty zsh autocompletion folder
    When I run "git-town completions zsh"
    Then it prints:
      """
      Zsh autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with an existing Git autocompletion file
    Given my computer has an existing Git autocompletion file for zsh
    When I run "git-town completions zsh"
    Then it prints the error:
      """
      Zsh autocompletions already exists
      """
    And my computer still has the original Git autocompletion file
