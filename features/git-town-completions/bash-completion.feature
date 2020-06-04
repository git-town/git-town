Feature: Installing Bash autocomplete definitions

  As a Bash shell user
  I want to use Git Town with autocompletion
  So that I can use this tool productively.


  Scenario: without existing Bash autocompletion folder
    Given my computer has no Bash autocompletion file
    When I run "git-town completions bash"
    Then it prints:
      """
      Bash autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with empty Bash autocompletion folder
    Given my computer has an empty Bash autocompletion folder
    When I run "git-town completions bash"
    Then it prints:
      """
      Bash autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with an existing Git autocompletion file
    Given my computer has an existing Git autocompletion file for Bash
    When I run "git-town completions bash"
    Then it prints the error:
      """
      Bash autocompletions already exists
      """
    And my computer still has the original Git autocompletion file
