@skipWindows
Feature: GitHub support

  Background:
    Given a Git repo with origin
    And a proposal for this branch does not exist
    And tool "open" is installed

  Scenario Outline: normal origin
    Given the origin is "<ORIGIN>"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      |         | Looking for proposal online ... ok                                 |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1 |

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
    Given the origin is "<ORIGIN>"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                       |
      | feature | git fetch --prune --tags                                                      |
      |         | Looking for proposal online ... ok                                            |
      |         | open https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |

    Examples:
      | ORIGIN                                              |
      | http://github.com/git-town/git-town.github.com.git  |
      | http://github.com/git-town/git-town.github.com      |
      | https://github.com/git-town/git-town.github.com.git |
      | https://github.com/git-town/git-town.github.com     |
      | git@github.com:git-town/git-town.github.com.git     |
      | git@github.com:git-town/git-town.github.com         |

  Scenario: URL-encodes hashtag
    Given the origin is "https://github.com/git-town/git-town"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | fix-#2 | feature | main   | local, origin |
    And the current branch is "fix-#2"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                             |
      | fix-#2 | git fetch --prune --tags                                            |
      |        | Looking for proposal online ... ok                                  |
      |        | open https://github.com/git-town/git-town/compare/fix-%232?expand=1 |

  Scenario: URL-encodes forward slashes
    Given the origin is "https://github.com/git-town/git-town"
    And the branches
      | NAME         | TYPE    | PARENT | LOCATIONS     |
      | test/feature | feature | main   | local, origin |
    And the current branch is "test/feature"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                                   |
      | test/feature | git fetch --prune --tags                                                  |
      |              | Looking for proposal online ... ok                                        |
      |              | open https://github.com/git-town/git-town/compare/test%2Ffeature?expand=1 |

  Scenario: stacked change with known parent
    Given the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                   |
      | child  | git fetch --prune --tags                                                  |
      |        | git checkout parent                                                       |
      | parent | git checkout child                                                        |
      |        | Looking for proposal online ... ok                                        |
      | child  | open https://github.com/git-town/git-town/compare/parent...child?expand=1 |
