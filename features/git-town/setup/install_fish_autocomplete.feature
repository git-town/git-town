Feature: Installing Fish Shell autocomplete definitions

  As a Fish shell user
  I want to be able to install the autocomplete definitions for Git Town with an easy command
  So that I can use this tool productively despite not having time for long installation procedures.


  Scenario: with a standard Fish setup
    When I run `git town install-fish-autocompletion`
    Then it runs the following shell commands to install the Fish shell autocompletion
      | COMMAND                                                                                                                       |
      | mkdir -p ~/.config/fish/completions                                                                                           |
      | curl -o ~/.config/fish/completions/git.fish https://raw.githubusercontent.com/Originate/git-town/master/autocomplete/git.fish |
