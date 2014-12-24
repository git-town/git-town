Feature: Git Town performs correctly depending on the environment

  Scenario Outline: Git Town commands being run not in a git repository
    Given I'm currently not in a git repository
    When I run `<command>` while allowing errors
    Then I see "This is not a git repository."

  Examples:
    | command             |
    |  git extract        |
    |  git hack           |
    |  git kill           |
    |  git pr             |
    |  git prune-branches |
    |  git repo           |
    |  git ship           |
    |  git sync           |
    |  git sync-fork      |


  Scenario Outline: `git town` commands being run not in a git repository
    Given I'm currently not in a git repository
    When I run `<command>` while allowing errors
    Then I <do_or_dont> see "This is not a git repository."

  Examples:
    | command                         | do_or_dont |
    |   git town                      |   don't    |
    |   git town help                 |   don't    |
    |   git town version              |   don't    |
    |   git town config               |   do       |
    |   git town main-branch          |   do       |
    |   git town non-feature-branches |   do       |
