@skipWindows
Feature: GitHub support

  Background:
    Given a Git repo clone
    Given tool "open" is installed

  Scenario Outline: normal origin
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given the current branch is "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """

    Examples:
      | ORIGIN                                     |
      | http://github.com/git-town/git-town.git    |
      | http://github.com/git-town/git-town        |
      | https://github.com/git-town/git-town.git   |
      | https://github.com/git-town/git-town       |
      | git@github.com:git-town/git-town.git       |
      | git@github.com:git-town/git-town           |
      | ssh://git@github.com/git-town/git-town.git |
      | ssh://git@github.com/git-town/git-town     |

  Scenario Outline: origin contains path that looks like a URL
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given the current branch is "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """

    Examples:
      | ORIGIN                                              |
      | http://github.com/git-town/git-town.github.com.git  |
      | http://github.com/git-town/git-town.github.com      |
      | https://github.com/git-town/git-town.github.com.git |
      | https://github.com/git-town/git-town.github.com     |
      | git@github.com:git-town/git-town.github.com.git     |
      | git@github.com:git-town/git-town.github.com         |

  Scenario Outline: proper URL encoding
    Given the branches
      | NAME          | TYPE    | PARENT | LOCATIONS     |
      | <BRANCH_NAME> | feature | main   | local, origin |
    Given the current branch is "<BRANCH_NAME>"
    And the origin is "https://github.com/git-town/git-town"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      <URL>
      """

    Examples:
      | BRANCH_NAME    | URL                                                                  |
      | feature-branch | https://github.com/git-town/git-town/compare/feature-branch?expand=1 |
      | feature_branch | https://github.com/git-town/git-town/compare/feature_branch?expand=1 |
      | fix-#2         | https://github.com/git-town/git-town/compare/fix-%232?expand=1       |
      | test/feature   | https://github.com/git-town/git-town/compare/test%2Ffeature?expand=1 |

  Scenario: stacked change with known parent
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the origin is "git@github.com:git-town/git-town.git"
    And the current branch is "child"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent...child?expand=1
      """
