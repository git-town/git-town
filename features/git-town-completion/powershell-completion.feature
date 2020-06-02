Feature: Installing Powershell Shell autocomplete definitions

  As a Powershell shell user
  I want to be able to install the autocomplete definitions for Git Town with an easy command
  So that I can use this tool productively despite not having time for long installation procedures.


  Scenario: without existing powershell autocompletion folder
    Given my computer has no powershell autocompletion file
    When I run "git-town completion powershell"
    Then it prints:
      """
      Git autocompletion for Powershell shell installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with empty powershell autocompletion folder
    Given my computer has an empty powershell autocompletion folder
    When I run "git-town completion powershell"
    Then it prints:
      """
      Git autocompletion for Powershell shell installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with an existing Git autocompletion file
    Given my computer has an existing Git autocompletion file
    When I run "git-town completion powershell"
    Then it prints the error:
      """
      Git autocompletion for Powershell shell already exists
      """
    And my computer still has the original Git autocompletion file
