Feature: support Git configuration that includes other files

  Scenario: global config file contains an include directive
    Given a Git repo with origin
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git push --tags          |
    Given the home directory contains file ".gitconfig" with content
      """
      [include]
      path = .gitconfig.user
      """
    And the home directory contains file ".gitconfig.user" with content
      """
      [user]
      name = User Name
      email = user@example.com
      """
