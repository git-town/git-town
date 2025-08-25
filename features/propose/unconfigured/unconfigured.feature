@messyoutput
Feature: ask for missing configuration

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And the origin is "https://github.com/git-town/git-town.git"
    When I run "git-town propose" and enter into the dialog:
      | DIALOG                      | KEYS            |
      | welcome                     | enter           |
      | aliases                     | enter           |
      | main branch                 | enter           |
      | perennial branches          |                 |
      | perennial regex             | enter           |
      | feature regex               | enter           |
      | contribution regex          | enter           |
      | observed regex              | enter           |
      | new branch type             | enter           |
      | unknown branch type         | enter           |
      | origin hostname             | enter           |
      | forge type                  | enter           |
      | github connector type       | enter           |
      | github token                | t o k e n enter |
      | token scope                 | enter           |
      | sync feature strategy       | enter           |
      | sync perennial strategy     | enter           |
      | sync prototype strategy     | enter           |
      | sync upstream               | enter           |
      | sync tags                   | enter           |
      | detached                    | enter           |
      | stash                       | enter           |
      | share new branches          | enter           |
      | push hook                   | enter           |
      | ship strategy               | enter           |
      | ship delete tracking branch | enter           |
      | config storage              | enter           |

  Scenario: result
    Then the main branch is now "main"
    And Git Town prints the error:
      """
      cannot propose the main branch
      """
