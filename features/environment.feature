Feature: Git Town performs correctly depending on the environment

  Scenario Outline: Git Town commands run outside of a Git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>`
    Then I get the error "This is not a git repository."

    Examples:
      | COMMAND               |
      | gt config             |
      | gt main-branch        |
      | gt perennial-branches |
      | gt hack feature       |
      | gt kill               |
      | gt new-pull-request   |
      | gt prune-branches     |
      | gt repo               |
      | gt ship               |
      | gt sync               |
