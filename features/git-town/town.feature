Feature: Show correct git town usage

  Scenario: invalid git town command
    When I run `git town invalidcommand`
    Then I see
      """
      'invalidcommand' is not a valid Git Town command

        Usage
        git town
        git town config
        git town help
        git town main-branch [<branchname>]
        git town non-feature-branches [(--add | --remove) <branchname>]
        git town version
      """
