Feature: print nice error message for invalid config file

  Scenario: Config file with invalid TOML content
    Given a Git repo clone
    And the configuration file:
      """
      push-new-branches =
      """
    When I run "git-town config"
    Then it prints the error:
      """
      the configuration file ".git-branches.yml" does not contain TOML-formatted content
      """
