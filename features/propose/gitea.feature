@skipWindows
Feature: Gitea support

  Background:
    Given a Git repo clone
    And tool "open" is installed

  Scenario Outline: normal origin
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://gitea.com/git-town/git-town/compare/main...feature
      """

    Examples:
      | ORIGIN                                    |
      | http://gitea.com/git-town/git-town.git    |
      | http://gitea.com/git-town/git-town        |
      | https://gitea.com/git-town/git-town.git   |
      | https://gitea.com/git-town/git-town       |
      | git@gitea.com:git-town/git-town.git       |
      | git@gitea.com:git-town/git-town           |
      | ssh://git@gitea.com/git-town/git-town.git |
      | ssh://git@gitea.com/git-town/git-town     |

  Scenario Outline: origin contains path that looks like a URL
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://gitea.com/git-town/git-town.gitea.com/compare/main...feature
      """

    Examples:
      | ORIGIN                                            |
      | http://gitea.com/git-town/git-town.gitea.com.git  |
      | http://gitea.com/git-town/git-town.gitea.com      |
      | https://gitea.com/git-town/git-town.gitea.com.git |
      | https://gitea.com/git-town/git-town.gitea.com     |
      | git@gitea.com:git-town/git-town.gitea.com.git     |
      | git@gitea.com:git-town/git-town.gitea.com         |

  Scenario Outline: proper URL encoding
    Given the branches
      | NAME          | TYPE    | PARENT | LOCATIONS     |
      | <BRANCH_NAME> | feature | main   | local, origin |
    And the current branch is "<BRANCH_NAME>"
    And the origin is "https://gitea.com/git-town/git-town"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      <URL>
      """

    Examples:
      | BRANCH_NAME    | URL                                                               |
      | feature-branch | https://gitea.com/git-town/git-town/compare/main...feature-branch |
      | feature_branch | https://gitea.com/git-town/git-town/compare/main...feature_branch |
      | fix-#2         | https://gitea.com/git-town/git-town/compare/main...fix-%232       |
      | test/feature   | https://gitea.com/git-town/git-town/compare/main...test%2Ffeature |

  Scenario: stacked change with known parent
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the origin is "git@gitea.com:git-town/git-town.git"
    And the current branch is "child"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://gitea.com/git-town/git-town/compare/parent...child
      """
