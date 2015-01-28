Feature: Initial configuration

  Scenario Outline: Running Git Town commands while Git Town is unconfigured
    Given I haven't configured Git Town yet
    When I run `<COMMAND>` while allowing errors
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


  Scenario: Enter invalid main branch
    Given I haven't configured Git Town yet
    When I run `git town config --setup` and enter "invalidbranch"
    Then I see
      """
      Please enter the name of the main dev branch (typically 'master' or 'development'):

        Error
        There is no branch named 'invalidbranch'

      """
    And Git Town is still not configured for this repository


  Scenario: Enter invalid non-feature branch
    Given I haven't configured Git Town yet
    And I have a branch named "master"
    When I run `git town config --setup` and enter "master" and "invalidbranch"
    Then I see
      """
      Please enter the name of the main dev branch (typically 'master' or 'development'):

      Git Town supports non-feature branches like 'release' or 'production'.
      These branches cannot be shipped and do not merge master when syncing.
      Please enter the names of all your non-feature branches as a comma separated list.
      Example: 'qa, production'

        Error
        There is no branch named 'invalidbranch'


      """
    And the main branch name is now configured as "master"
    And my non-feature branches are still not configured
