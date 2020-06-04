Feature: Installing Powershell autocomplete definitions

  As a Powershell user
  I want to use Git Town with autocompletion
  So that I can use this tool productively.


  Scenario: without existing PowerShell autocompletion folder
    Given my computer has no PowerShell autocompletion file
    When I run "git-town completions powershell"
    Then it prints:
      """
      PowerShell autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with empty PowerShell autocompletion folder
    Given my computer has an empty PowerShell autocompletion folder
    When I run "git-town completions powershell"
    Then it prints:
      """
      PowerShell autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with an existing Git autocompletion file
    Given my computer has an existing Git autocompletion file for PowerShell
    When I run "git-town completions powershell"
    Then it prints the error:
      """
      PowerShell autocompletions already exists
      """
    And my computer still has the original Git autocompletion file
