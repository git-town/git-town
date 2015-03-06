@modifies-fish-completions
Feature: Installing Fish Shell autocomplete definitions

  As a Fish shell user
  I want to be able to install the autocomplete definitions for Git Town with an easy command
  So that I can use this tool productively despite not having time for long installation procedures.


  Scenario: without existing fish autocompletion folder
    Given I have no fish autocompletion folder
    When I run `git town install-fish-autocompletion`
    Then it runs the following shell commands to install the Fish shell autocompletion
      | COMMAND                                                                                   |
      | mkdir -p ~/.config/fish/completions                                                       |
      | ln -s <%= GIT_TOWN_DIRECTORY %>/autocomplete/git.fish ~/.config/fish/completions/git.fish |


  Scenario: with empty fish autocompletion folder
    Given I have an empty fish autocompletion folder
    When I run `git town install-fish-autocompletion`
    Then it runs the following shell commands to install the Fish shell autocompletion
      | COMMAND                                                                                   |
      | ln -s <%= GIT_TOWN_DIRECTORY %>/autocomplete/git.fish ~/.config/fish/completions/git.fish |


  Scenario: with an existing Git autocompletion file
    Given I have an existing Git autocompletion file
    When I run `git town install-fish-autocompletion`
    Then it runs no shell commands
    And I get the error "Git autocompletion for Fish shell already exists"
    And I get the error "Operation aborted"
    And I still have my original Git autocompletion file


  Scenario: with existing Git autocompletion symlink
    Given I already have the Git autocompletion symlink
    When I run `git town install-fish-autocompletion`
    Then it runs no shell commands
    And I get the error "Git autocompletion for Fish shell already exists"
    And I get the error "Operation aborted"
    And I still have my original Git autocompletion symlink
