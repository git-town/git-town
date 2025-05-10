@skipWindows
Feature: Gitea support

  Background:
    Given a Git repo with origin
    And tool "open" is installed
    And a proposal for this branch does not exist

  Scenario Outline: normal origin
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the origin is "<ORIGIN>"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                         |
      | feature | git fetch --prune --tags                                        |
      | (none)  | Looking for proposal online ... ok                              |
      |         | open https://gitea.com/git-town/git-town/compare/main...feature |

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
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                   |
      | feature | git fetch --prune --tags                                                  |
      | (none)  | Looking for proposal online ... ok                                        |
      |         | open https://gitea.com/git-town/git-town.gitea.com/compare/main...feature |

    Examples:
      | ORIGIN                                            |
      | http://gitea.com/git-town/git-town.gitea.com.git  |
      | http://gitea.com/git-town/git-town.gitea.com      |
      | https://gitea.com/git-town/git-town.gitea.com.git |
      | https://gitea.com/git-town/git-town.gitea.com     |
      | git@gitea.com:git-town/git-town.gitea.com.git     |
      | git@gitea.com:git-town/git-town.gitea.com         |

  Scenario: URL-encodes hashtag
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | fix-#2 | feature | main   | local, origin |
    And the current branch is "fix-#2"
    And the origin is "https://gitea.com/git-town/git-town"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                          |
      | fix-#2 | git fetch --prune --tags                                         |
      | (none) | Looking for proposal online ... ok                               |
      |        | open https://gitea.com/git-town/git-town/compare/main...fix-%232 |

  Scenario: URL-encodes forward slashes
    Given the branches
      | NAME         | TYPE    | PARENT | LOCATIONS     |
      | test/feature | feature | main   | local, origin |
    And the current branch is "test/feature"
    And the origin is "https://gitea.com/git-town/git-town"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                                |
      | test/feature | git fetch --prune --tags                                               |
      | (none)       | Looking for proposal online ... ok                                     |
      |              | open https://gitea.com/git-town/git-town/compare/main...test%2Ffeature |

  Scenario: stacked change with known parent
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the origin is "git@gitea.com:git-town/git-town.git"
    And the current branch is "child"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                         |
      | child  | git fetch --prune --tags                                        |
      | (none) | Looking for proposal online ... ok                              |
      | child  | git checkout parent                                             |
      | parent | git checkout child                                              |
      | (none) | open https://gitea.com/git-town/git-town/compare/parent...child |
