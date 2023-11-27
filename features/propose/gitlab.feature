@skipWindows
Feature: GitLab support

  Background:
    Given tool "open" is installed

  Scenario Outline: creating proposals
    Given the current branch is a feature branch "feature"
    And the origin is "<REPO ORIGIN>"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      <BROWSER URL>
      """

    Examples:
      | REPO ORIGIN                                  | BROWSER URL                                                                                                                                      |
      | https://gitlab.com/kadu/kadu.git             | https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main                 |
      | git@gitlab.com:kadu/kadu.git                 | https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main                 |
      | git@gitlab.com:gitlab-com/www-gitlab-com.git | https://gitlab.com/gitlab-com/www-gitlab-com/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |

  Scenario: nested feature branch with known parent
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the origin is "git@gitlab.com:kadu/kadu.git"
    And the current branch is "child"
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://gitlab.com/kadu/kadu/-/merge_requests/new?merge_request%5Bsource_branch%5D=child&merge_request%5Btarget_branch%5D=parent
      """
