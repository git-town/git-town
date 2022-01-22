@skipWindows
Feature: "git town config setup" on completely empty repo

  To make the configuration dialog intelligent
  I don't want to be asked for perennial branches if there are no branches that could be perennial.

  Background:
    Given I haven't configured Git Town yet
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER        |
      | Please specify the main development branch | [DOWN][ENTER] |

  Scenario: result
    Then the main branch is now configured as "main"
    And my repo is now configured with no perennial branches
