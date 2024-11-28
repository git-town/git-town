@messyoutput
Feature: switch branches from detached head

  Background:
    Given a Git repo with origin
    And I ran "git checkout HEAD^"
    And inspect the repo
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    And inspect the repo

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout beta |
    And the current branch is now "beta"
