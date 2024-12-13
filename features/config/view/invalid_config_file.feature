Feature: print nice error message for invalid config file

  Scenario: Config file with invalid TOML content
    Given a Git repo with origin
    And the configuration file:
      """
      wrong =
      """
    When I run "git-town config"
    Then Git Town prints the error:
      """
      the configuration file ".git-branches.yml" does not contain TOML-formatted content
      """
