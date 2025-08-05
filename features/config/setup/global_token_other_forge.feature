@messyoutput
Feature: a global API token of another forge exists

  Scenario: on GitHub, with global GitLab token
    Given a Git repo with origin
    And my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And global Git setting "git-town.gitlab-token" is "987654"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS            | DESCRIPTION                                 |
      | welcome                     | enter           |                                             |
      | aliases                     | enter           |                                             |
      | main branch                 | enter           |                                             |
      | perennial branches          |                 | no input here since the dialog doesn't show |
      | perennial regex             | enter           |                                             |
      | feature regex               | enter           |                                             |
      | contribution regex          | enter           |                                             |
      | observed regex              | enter           |                                             |
      | new branch type             | enter           |                                             |
      | unknown branch type         | enter           |                                             |
      | dev remote                  | enter           |                                             |
      | origin hostname             | enter           |                                             |
      | forge type                  | enter           |                                             |
      | github connector type       | enter           |                                             |
      | github token                | g h t o k enter |                                             |
      | token scope                 | enter           |                                             |
      | sync feature strategy       | enter           |                                             |
      | sync perennial strategy     | enter           |                                             |
      | sync prototype strategy     | enter           |                                             |
      | sync upstream               | enter           |                                             |
      | sync tags                   | enter           |                                             |
      | share new branches          | enter           |                                             |
      | push hook                   | enter           |                                             |
      | ship strategy               | enter           |                                             |
      | ship delete tracking branch | enter           |                                             |
      | config storage              | enter           |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config git-town.github-token ghtok   |
      | git config git-town.github-connector api |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" is now "ghtok"
    And local Git setting "git-town.gitlab-token" now doesn't exist
    And global Git setting "git-town.gitlab-token" is still "987654"
