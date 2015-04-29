Feature: Git Town performs correctly depending on the environment

  Scenario Outline: Git Town commands run outside of a Git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>`
    Then I get the error
      """
      error: This is not a git repository.
      """

    Examples:
      | COMMAND                       |
      | git town config               |
      | git town main-branch          |
      | git town non-feature-branches |
      | git create-pull-request       |
      | git extract                   |
      | git hack                      |
      | git kill                      |
      | git prune-branches            |
      | git repo                      |
      | git ship                      |
      | git sync                      |
      | git sync-fork                 |
