Feature: support Git configuration that includes other files

  @debug @this
  Scenario: global config file contains an include directive
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
    When I run "git-town sync"
    Then it prints:
      """
      xxx
      """
