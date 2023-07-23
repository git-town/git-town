Feature: display debug statistics

  Scenario: Git Town command ran successfully
    Given I ran "git-town sync"
    When I run "git-town status --debug"
    Then it runs the debug commands
      | git config -lz --local        |
      | git config -lz --global       |
      | git rev-parse                 |
      | git rev-parse --show-toplevel |
      | git version                   |
