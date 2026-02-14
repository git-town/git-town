@messyoutput
Feature: Configure a different development remote

  Background:
    Given a Git repo with origin
    And an additional "fork" remote with URL "https://github.com/forked/repo"
    When I run "git-town init" and enter into the dialogs:
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
      #:schema https://raw.githubusercontent.com/git-town/git-town/refs/heads/main/docs/git-town.schema.json

      # See https://www.git-town.com/configuration-file for details

      [branches]
      main = "main"

      [hosting]
      dev-remote = "fork"
      """
