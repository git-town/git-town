@messyoutput
Feature: migrate connector configuration from Git metadata to config file

  Background:
    Given a Git repo with origin
    And the main branch is "main"
    And local Git setting "git-town.forge-type" is "github"
    And local Git setting "git-town.github-connector-type" is "api"
    And local Git setting "git-town.hosting-origin-hostname" is "github.example.com"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                | KEYS       |
      | welcome               | enter      |
      | aliases               | enter      |
      | main branch           | enter      |
      | perennial branches    | enter      |
      | dev-remote            | enter      |
      | github connector type | down enter |
      | enter all             | down enter |
      | config storage        | down enter |

  Scenario: result
    Then local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-connector-type" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And the configuration file is now:
      """
      # See https://www.git-town.com/configuration-file for details

      [branches]
      main = "main"

      [hosting]
      forge-type = "github"
      github-connector-type = "gh"
      origin-hostname = "github.example.com"

      [sync]
      prototype-strategy = "merge"
      """

  Scenario: undo
    When I run "git-town undo"
    Then local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.github-connector-type" is now "api"
    And local Git setting "git-town.hosting-origin-hostname" is now "github.example.com"
    And the main branch is now "main"
