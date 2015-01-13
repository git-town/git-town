Feature: Git Town performs correctly depending on the environment

  Scenario Outline: Git Town commands run outside of a Git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>` while allowing errors
    Then I <DO_OR_DONT> see "This is not a git repository."

  Examples:
    | COMMAND                       | DO_OR_DONT |
    | git town                      | don't      |
    | git town help                 | don't      |
    | git town version              | don't      |
    | git town config               | do         |
    | git town main-branch          | do         |
    | git town non-feature-branches | do         |
    | git extract                   | do         |
    | git hack                      | do         |
    | git kill                      | do         |
    | git pr                        | do         |
    | git prune-branches            | do         |
    | git repo                      | do         |
    | git ship                      | do         |
    | git sync                      | do         |
    | git sync-fork                 | do         |
