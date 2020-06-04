Feature: Installing Fish Shell autocomplete definitions

  As a Fish shell user
  I want to use Git Town with autocompletion
  So that I can use this tool productively.


  Scenario: without existing fish autocompletion folder
    Given my computer has no fish autocompletion file
    When I run "git-town completions fish"
    Then it prints:
      """
      Fish autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with empty fish autocompletion folder
    Given my computer has an empty fish autocompletion folder
    When I run "git-town completions fish"
    Then it prints:
      """
      Fish autocompletions installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with an existing Git autocompletion file
    Given my computer has an existing Git autocompletion file for fish
    When I run "git-town completions fish"
    Then it prints the error:
      """
      Fish autocompletions already exists
      """
    And my computer still has the original Git autocompletion file
