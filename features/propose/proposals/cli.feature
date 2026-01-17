@skipWindows
Feature: propose with embedded lineage

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And tool "open" is installed
    When I run "git-town propose"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                        |
      | feature | git fetch --prune --tags                                                                       |
      |         | Finding proposal from feature into main ... none                                               |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1&title=my_title&body=my_body |
