Feature: Git Town performs correctly depending on the environment

  Scenario Outline: Git Town commands run outside of a Git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>`
    Then I get the error "This is not a Git repository."

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
