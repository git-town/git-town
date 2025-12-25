@messyoutput
Feature: enter Azure DevOps

  Background:
    Given a Git repo with origin

  Scenario: auto-detected Azure Devops platform
    And my repo's "origin" remote is "git@ssh.dev.azure.com:v3/kevingoslar/tikibase/tikibase"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS  | DESCRIPTION                                 |
      | welcome            | enter |                                             |
      | aliases            | enter |                                             |
      | main branch        | enter |                                             |
      | perennial branches |       | no input here since the dialog doesn't show |
      | origin hostname    | enter |                                             |
      | forge type         | enter | auto-detect                                 |
      | enter all          | enter |                                             |
      | config storage     | enter | git metadata                                |
    Then Git Town runs no commands
    And local Git setting "git-town.forge-type" still doesn't exist

  Scenario: select Azure DevOps manually
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS       | DESCRIPTION                                 |
      | welcome            | enter      |                                             |
      | aliases            | enter      |                                             |
      | main branch        | enter      |                                             |
      | perennial branches |            | no input here since the dialog doesn't show |
      | origin hostname    | enter      |                                             |
      | forge type         | down enter |                                             |
      | enter all          | enter      |                                             |
      | config storage     | enter      | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                    |
      | git config git-town.forge-type azuredevops |
    And local Git setting "git-town.forge-type" is now "azuredevops"
