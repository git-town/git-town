@messyoutput
Feature: migrate GitLab connector configuration from Git metadata to config file

  Background:
    Given a Git repo with origin
    And the main branch is "main"
    And local Git setting "git-town.forge-type" is "gitlab"
    And local Git setting "git-town.gitlab-connector-type" is "api"
    And local Git setting "git-town.hosting-origin-hostname" is "gitlab.example.com"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                | KEYS       |
      | welcome               | enter      |
      | aliases               | enter      |
      | main branch           | enter      |
      | perennial branches    | enter      |
      | dev-remote            | enter      |
      | gitlab connector type | down enter |
      | enter all             | down enter |
      | config storage        | down enter |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                             |
      | git config --unset git-town.forge-type              |
      | git config --unset git-town.gitlab-connector-type   |
      | git config --unset git-town.hosting-origin-hostname |
      | git config --unset git-town.main-branch             |
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.gitlab-connector-type" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And the configuration file is now:
      """
      # See https://www.git-town.com/configuration-file for details

      [branches]
      main = "main"

      [hosting]
      forge-type = "gitlab"
      gitlab-connector-type = "glab"
      origin-hostname = "gitlab.example.com"

      [sync]
      prototype-strategy = "merge"
      """
    And the main branch is now not set

  Scenario: undo
    When I run "git-town undo"
    Then local Git setting "git-town.forge-type" is now "gitlab"
    And local Git setting "git-town.gitlab-connector-type" is now "api"
    And local Git setting "git-town.hosting-origin-hostname" is now "gitlab.example.com"
    And the main branch is now "main"
