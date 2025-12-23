@skipWindows
Feature: GitLab support

  Background:
    Given a Git repo with origin
    And a proposal for this branch does not exist
    And tool "open" is installed

  Scenario Outline: creating proposals
    Given the origin is "<REPO ORIGIN>"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | Looking for proposal online ... ok |
      |         | open <BROWSER URL>                 |

    Examples:
      | REPO ORIGIN                                  | BROWSER URL                                                                                                                                      |
      | https://gitlab.com/kadu/kadu.git             | https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main                 |
      | git@gitlab.com:kadu/kadu.git                 | https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main                 |
      | git@gitlab.com:gitlab-com/www-gitlab-com.git | https://gitlab.com/gitlab-com/www-gitlab-com/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |

  Scenario: stacked change with known parent
    Given the origin is "git@gitlab.com:kadu/kadu.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                                               |
      | child  | git fetch --prune --tags                                                                                                              |
      |        | git checkout parent                                                                                                                   |
      | parent | git checkout child                                                                                                                    |
      |        | Looking for proposal online ... ok                                                                                                    |
      | child  | open https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=child&merge_request%5Btarget_branch%5D=parent |
