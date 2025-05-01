Feature: the Git configuration uses includeIf

  Scenario: global Git config uses includeIf
    Given a Git repo with origin
    And the home directory contains file ".gitconfig" with content
      """
      [includeIf "hasconfig:remote.*.url:git@github.com*/**"]
      path = .gitconfig-personal
      """
    And the home directory contains file ".gitconfig-personal" with content
      """
      [user]
        name = The User
        email = user@acme.com
      """
    When I run "git-town hack new"
    Then Git Town prints the error:
      """
      please set the Git user email by running: git config --global user.email "<your email>"
      """
