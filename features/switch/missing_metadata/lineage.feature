@messyoutput
Feature: switch branches that have no lineage information

  Scenario: repo contains a branch without known parent
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE   | LOCATIONS |
      | alpha | (none) | local     |
      | beta  | (none) | local     |
    And the current branch is "alpha"
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout beta |
    And the current branch is now "beta"
