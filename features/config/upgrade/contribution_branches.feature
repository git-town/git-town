Feature: Automatically remove obsolete branch lists setting

  Scenario Outline:
    Given a local Git repo
    And local Git setting "git-town.<BRANCHTYPE>-branches" is "one two"
    When I run "git-town config"
    Then Git Town prints:
      """
      Inlining deprecated branch list "git-town.<BRANCHTYPE>-branches"
      """
    And local Git setting "git-town.<BRANCHTYPE>-branches" now doesn't exist
    And local Git setting "git-town-branch.one.branchtype" is now "<BRANCHTYPE>"
    And local Git setting "git-town-branch.two.branchtype" is now "<BRANCHTYPE>"

    Examples:
      | BRANCHTYPE   |
      | contribution |
      | observed     |
      | parked       |
      | prototype    |
