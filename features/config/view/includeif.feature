Feature: the Git configuration uses includeIf

  @this
  Scenario: global Git config uses includeIf
    Given a Git repo with origin
    And Git Town is not configured
    And file ".gitconfig" in the homedirectory
      """
      [includeIf "hasconfig:remote.*.url:git@github.com*/**"]
      path = .gitconfig-personal
      """
    And file ".gitconfig-personal" in the homedirectory
      """
      [user]
        name = The User
        email = user@acme.com
      """
    When I run "git-town config"
    Then Git Town prints:
      """
      xxx
      """
