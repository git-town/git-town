@skipWindows
Feature: multi-platform support

  Background:
    Given a Git repo with origin

  Scenario Outline: supported tool installed
    Given the origin is "https://github.com/git-town/git-town.git"
    And tool "<TOOL>" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                     |
      | main   | frontend | <TOOL> https://github.com/git-town/git-town |

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |
