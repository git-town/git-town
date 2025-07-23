@messyoutput
Feature: ask for missing configuration information

  Scenario: run unconfigured
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                      | KEYS  |
      | welcome                     | enter |
      | aliases                     | enter |
      | main branch                 | enter |
      | perennial branches          |       |
      | perennial regex             | enter |
      | feature regex               | enter |
      | contribution regex          | enter |
      | observed regex              | enter |
      | new branch type             | enter |
      | unknown branch type         | enter |
      | origin hostname             | enter |
      | forge type                  | enter |
      | sync feature strategy       | enter |
      | sync perennial strategy     | enter |
      | sync prototype strategy     | enter |
      | sync upstream               | enter |
      | sync tags                   | enter |
      | share new branches          | enter |
      | push hook                   | enter |
      | ship strategy               | enter |
      | ship delete tracking branch | enter |
      | config storage              | enter |
    And the main branch is now "main"
