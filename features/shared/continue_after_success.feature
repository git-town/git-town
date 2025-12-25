Feature: continue after successful command

  Scenario Outline:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE |
      | feature | local, origin | commit  |
    And local Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "feature"
    And tool "open" is installed
    And I ran "git-town <COMMAND>"
    When I run "git-town continue"
    Then Git Town prints:
      """
      nothing to continue
      """

    Examples:
      | COMMAND              |
      |                      |
      | append new           |
      | completions fish     |
      | config               |
      | diff-parent          |
      | hack new             |
      | help                 |
      | delete feature       |
      | offline              |
      | prepend new          |
      | propose              |
      | rename foo           |
      | repo                 |
      | ship feature -m done |
      | sync                 |
      | --version            |
