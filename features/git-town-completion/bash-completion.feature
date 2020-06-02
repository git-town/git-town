Feature: Installing Bash Shell autocomplete definitions

  As a Bash shell user
  I want to be able to install the autocomplete definitions for Git Town with an easy command
  So that I can use this tool productively despite not having time for long installation procedures.


  Scenario: without existing bash autocompletion folder
    Given my computer has no bash autocompletion file
    When I run "git-town completion bash"
    Then it prints:
      """
      Git autocompletion for Bash shell installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with empty bash autocompletion folder
    Given my computer has an empty bash autocompletion folder
    When I run "git-town completion bash"
    Then it prints:
      """
      Git autocompletion for Bash shell installed
      """
    And my computer now has a Git autocompletion file


  Scenario: with an existing Git autocompletion file
    Given my computer has an existing Git autocompletion file
    When I run "git-town completion bash"
    Then it prints the error:
      """
      Git autocompletion for Bash shell already exists
      """
    And my computer still has the original Git autocompletion file
