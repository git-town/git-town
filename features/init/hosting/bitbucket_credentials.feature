@messyoutput
Feature: enter the Bitbucket credentials

  Background:
    Given a Git repo with origin

  Scenario: auto-detected Bitbucket Cloud platform
    And my repo's "origin" remote is "git@bitbucket.org:git-town/git-town.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG              | KEYS                | DESCRIPTION                                 |
      | welcome             | enter               |                                             |
      | aliases             | enter               |                                             |
      | main branch         | enter               |                                             |
      | perennial branches  |                     | no input here since the dialog doesn't show |
      | origin hostname     | enter               |                                             |
      | forge type          | enter               | auto-detect                                 |
      | bitbucket username  | u s e r enter       |                                             |
      | bitbucket api token | a p i - t o k enter |                                             |
      | token scope         | enter               |                                             |
      | enter all           | enter               |                                             |
      | config storage      | enter               | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config git-town.bitbucket-username user     |
      | git config git-town.bitbucket-api-token api-tok |
    And local Git setting "git-town.bitbucket-api-token" is now "api-tok"
    And local Git setting "git-town.bitbucket-username" is now "user"
    And local Git setting "git-town.forge-type" still doesn't exist

  Scenario: select Bitbucket Data Center manually
    When I run "git-town init" and enter into the dialog:
      | DIALOG              | KEYS                  | DESCRIPTION                                 |
      | welcome             | enter                 |                                             |
      | aliases             | enter                 |                                             |
      | main branch         | enter                 |                                             |
      | perennial branches  |                       | no input here since the dialog doesn't show |
      | origin hostname     | enter                 |                                             |
      | forge type          | down down down enter  |                                             |
      | bitbucket username  | u s e r enter         |                                             |
      | bitbucket api token | h t t p - t o k enter |                                             |
      | token scope         | enter                 |                                             |
      | enter all           | enter                 |                                             |
      | config storage      | enter                 | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                             |
      | git config git-town.bitbucket-username user         |
      | git config git-town.bitbucket-api-token http-tok    |
      | git config git-town.forge-type bitbucket-datacenter |
    And local Git setting "git-town.bitbucket-api-token" is now "http-tok"
    And local Git setting "git-town.bitbucket-username" is now "user"
    And local Git setting "git-town.forge-type" is now "bitbucket-datacenter"
