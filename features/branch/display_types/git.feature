Feature: specify which branch types are displayed via Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | alpha        | feature      | main   | local, origin |
      | beta         | feature      | alpha  | local, origin |
      | gamma        | feature      | beta   | local, origin |
      | observed     | observed     |        | local, origin |
      | contribution | contribution |        | local, origin |
      | prototype    | prototype    | main   | local         |
      | parked       | parked       | main   | local         |
      | perennial    | perennial    |        | local, origin |
    And the current branch is "beta"

  Scenario: show all types
    Given Git setting "git-town.display-types" is "all"
    When I run "git-town branch"
    Then Git Town prints:
      """
        main  (main)
          alpha  (feature)
      *     beta  (feature)
              gamma  (feature)
          parked  (parked)
          prototype  (prototype)
        contribution  (contribution)
        observed  (observed)
        perennial  (perennial)
      """

  Scenario Outline: show only the given branch types
    Given Git setting "git-town.display-types" is "<VALUE>"
    When I run "git-town branch"
    Then Git Town prints:
      """
        main
          alpha
      *     beta
              gamma
          parked
          prototype  (prototype)
        contribution
        observed  (observed)
        perennial
      """

    Examples:
      | VALUE              |
      | prototype observed |
      | prototype+observed |
      | prototype-observed |
      | prototype_observed |
      | prototype&observed |

  Scenario: show no types
    Given Git setting "git-town.display-types" is "no"
    When I run "git-town branch"
    Then Git Town prints:
      """
        main
          alpha
      *     beta
              gamma
          parked
          prototype
        contribution
        observed
        perennial
      """

  Scenario Outline: show all except the given branch types
    Given Git setting "git-town.display-types" is "<VALUE>"
    When I run "git-town branch"
    Then Git Town prints:
      """
        main  (main)
          alpha  (feature)
      *     beta  (feature)
              gamma  (feature)
          parked  (parked)
          prototype
        contribution  (contribution)
        observed
        perennial  (perennial)
      """

    Examples:
      | VALUE                 |
      | no prototype observed |
      | no+prototype+observed |
      | no-prototype-observed |
      | no&prototype&observed |
      | no_prototype_observed |
      | no prototype-observed |
