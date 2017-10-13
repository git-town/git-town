Feature: Git Town performs correctly depending on the environment

  Scenario Outline: Git Town commands run outside of a Git repository
    Given my workspace is currently not a Git repository
    When I run `<COMMAND>`
    Then Git Town prints the error "This is not a Git repository"

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
