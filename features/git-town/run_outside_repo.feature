Feature: Require running inside a Git repository

  To learn to use Git Town correctly
  When running commands that require a repo outside of a repo
  I want to see guidance telling me that a repo is required.

  Scenario Outline:
    Given my workspace is currently not a Git repo
    When I run "<COMMAND>"
    Then it prints the error:
      """
      this is not a Git repository
      """

    Examples:
      | COMMAND                     |
      | git-town config             |
      | git-town main-branch        |
      | git-town perennial-branches |
      | git-town hack feature       |
      | git-town kill               |
      | git-town new-pull-request   |
      | git-town prune-branches     |
      | git-town repo               |
      | git-town ship               |
      | git-town sync               |
