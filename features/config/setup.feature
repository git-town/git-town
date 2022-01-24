@skipWindows
Feature: Initial configuration

  As a user who hasn't configured Git Town yet
  I want to have a simple, dedicated setup command
  So that I can configure it safely before using any Git Town command

  Scenario: succeeds on valid main branch and perennial branch names
    Given my repo has the feature branches "production" and "dev"
    And I haven't configured Git Town yet
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER                      |
      | Please specify the main development branch | [DOWN][ENTER]               |
      | Please specify perennial branches          | [SPACE][DOWN][SPACE][ENTER] |
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "dev" and "production"

  Scenario: does not prompt for perennial branches if there is only the main branch
    Given I haven't configured Git Town yet
    When I run "git-town config setup" and answer the prompts:
      | PROMPT                                     | ANSWER        |
      | Please specify the main development branch | [DOWN][ENTER] |
    Then the main branch is now configured as "main"
    And my repo is now configured with no perennial branches
