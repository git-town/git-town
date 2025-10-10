Feature: print nice error message for invalid alternative config file

  Scenario: Config file with alternative filename and invalid TOML content
    Given a Git repo with origin
    And an uncommitted file with name ".git-town.toml" and content:
      """
      wrong =
      """
    When I run "git-town config"
    Then Git Town prints the error:
      """
      the configuration file ".git-town.toml" does not contain TOML-formatted content
      """
