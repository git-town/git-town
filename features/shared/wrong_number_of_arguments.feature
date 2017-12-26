Feature: wrong number of arguments

  As a developer providing the wrong number of arguments
  I should be reminded of how many arguments the command expects
  So that I can use it correctly without having to look that fact up in the readme.


  Scenario Outline: <DESCRIPTION>
    When I run `git-town <CMD>`
    Then it runs no commands
    And it prints the error "Usage:"

    Examples:
      | CMD                            |
      | append                         |
      | append arg1 arg2               |
      | hack                           |
      | hack arg1 arg2                 |
      | new-branch-push-flag arg1 arg2 |
      | kill arg1 arg2                 |
      | main-branch arg1 arg2          |
      | new-pull-request arg1          |
      | perennial-branches arg1        |
      | prune-branches arg1            |
      | pull-branch-strategy arg1 arg2 |
      | rename-branch                  |
      | rename-branch arg1 arg2 arg3   |
      | repo arg1                      |
      | sync arg1                      |
      | ship arg1 arg2                 |
