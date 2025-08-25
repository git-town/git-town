@messyoutput
Feature: Configure a different development remote

  Background:
    Given a Git repo with origin
    And an additional "fork" remote with URL "https://github.com/forked/repo"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG             | KEYS       |
      | welcome            | enter      |
      | aliases            | enter      |
      | main branch        | enter      |
      | perennial branches |            |
      | dev-remote         | up enter   |
      | origin hostname    | enter      |
      | forge type         | enter      |
      | enter all          | enter      |
      | config storage     | down enter |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                 |
      | git config --unset git-town.main-branch |
    And the configuration file is now:
      """
      # More info around this file at https://www.git-town.com/configuration-file

      [branches]
      main = "main"

      [hosting]
      dev-remote = "fork"
      """
