@skipWindows
Feature: GitLab support

  Background:
    Given a Git repo with origin
    And tool "open" is installed
    And a proposal for this branch does not exist

  Scenario Outline: creating proposals
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the origin is "<REPO ORIGIN>"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      | <none>  | Looking for proposal online ... ok      |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      | <none>  | open <BROWSER URL>                      |

    Examples:
      | REPO ORIGIN                                  | BROWSER URL                                                                                                                                      |
      | https://gitlab.com/kadu/kadu.git             | https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main                 |
      | git@gitlab.com:kadu/kadu.git                 | https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main                 |
      | git@gitlab.com:gitlab-com/www-gitlab-com.git | https://gitlab.com/gitlab-com/www-gitlab-com/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |

  Scenario: stacked change with known parent
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the origin is "git@gitlab.com:kadu/kadu.git"
    And the current branch is "child"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                                               |
      | child  | git fetch --prune --tags                                                                                                              |
      | <none> | Looking for proposal online ... ok                                                                                                    |
      | child  | git checkout parent                                                                                                                   |
      | parent | git merge --no-edit --ff main                                                                                                         |
      |        | git merge --no-edit --ff origin/parent                                                                                                |
      |        | git checkout child                                                                                                                    |
      | child  | git merge --no-edit --ff parent                                                                                                       |
      |        | git merge --no-edit --ff origin/child                                                                                                 |
      | <none> | open https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=child&merge_request%5Btarget_branch%5D=parent |
