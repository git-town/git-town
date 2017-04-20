Feature: Installing Fish Shell autocomplete definitions

  As a Fish shell user
  I want to be able to install the autocomplete definitions for Git Town with an easy command
  So that I can use this tool productively despite not having time for long installation procedures.


  Scenario: without existing fish autocompletion folder
    Given I have no fish autocompletion folder
    When I run `gt install-fish-autocompletion`
    Then I see "Git autocompletion for Fish shell installed"
    And I have a Git autocompletion file


  Scenario: with empty fish autocompletion folder
    Given I have an empty fish autocompletion folder
    When I run `gt install-fish-autocompletion`
    Then I see "Git autocompletion for Fish shell installed"
    And I have a Git autocompletion file


  Scenario: with an existing Git autocompletion file
    Given I have an existing Git autocompletion file
    When I run `gt install-fish-autocompletion`
    Then I get the error "Git autocompletion for Fish shell already exists"
    And I still have my original Git autocompletion file


  Scenario: with existing Git autocompletion symlink
    Given I have an existing Git autocompletion symlink
    When I run `gt install-fish-autocompletion`
    Then I get the error "Git autocompletion for Fish shell already exists"
    And I still have my original Git autocompletion symlink
