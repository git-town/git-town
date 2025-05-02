Feature: the Git configuration uses includeIf

  Scenario: global Git config uses includeIf
    Given a Git repo with origin
    And the home directory contains file ".gitconfig" with content
      """
      [includeIf "onbranch:main"]
      path = .gitconfig-personal
      """
    And the home directory contains file ".gitconfig-personal" with content
      """
      [user]
        name = The User
        email = user@acme.com
      """
    When I run "git-town hack new"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git checkout -b new      |
