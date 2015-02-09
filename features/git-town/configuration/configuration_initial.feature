Feature: Initial configuration

  Background:
    Given I haven't configured Git Town yet


  Scenario Outline: Running Git Town commands while Git Town is unconfigured
    When I run `<COMMAND>`
    Then I see
      """
      Git Town hasn't been configured for this repository.
      Please run 'git town config --setup'.
      Would you like to do that now? y/n
      """

    Examples:
    | COMMAND            |
    | git extract        |
    | git hack           |
    | git kill           |
    | git pr             |
    | git prune-branches |
    | git repo           |
    | git ship           |
    | git sync           |
    | git sync-fork      |


  Scenario: Enter non-existent main branch
    When I run `git town config --setup` and enter "nonexistent"
    Then I get the error
      """
        Error
        There is no branch named 'nonexistent'
      """
    And Git Town is still not configured for this repository


  Scenario: Enter valid main branch and non-existent non-feature branch
    Given I have a branch named "master"
    When I run `git town config --setup` and enter "master" and "nonexistent"
    Then I get the error
      """
        Error
        There is no branch named 'nonexistent'
      """
    And the main branch name is now configured as "master"
    And my non-feature branches are still not configured


  Scenario: Enter valid main branch and valid non-feature branches
    Given I have branches named "dev" and "qa"
    When I run `git town config --setup` and enter "dev" and "qa"
    Then the main branch name is now configured as "dev"
    And my non-feature branches are now configured as "qa"


  Scenario: Enter valid main branch and invalid non-feature branches
    Given I have branches named "dev" and "qa"
    When I run `git town config --setup` and enter "dev" and "dev"
    Then I get the error
      """
        Error
        'dev' is already set as the main branch
      """
    And the main branch name is now configured as "dev"
    And my non-feature branches are still not configured
