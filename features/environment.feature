Feature: Git Town performs correctly depending on the environment

  Scenario Outline: Git Town commands run outside of a Git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>`, it errors
    Then I see "error: This is not a git repository."

    Examples:
      | COMMAND                       |
      | git town config               |
      | git town main-branch          |
      | git town non-feature-branches |
      | git extract                   |
      | git hack                      |
      | git kill                      |
      | git pr                        |
      | git prune-branches            |
      | git repo                      |
      | git ship                      |
      | git sync                      |
      | git sync-fork                 |
