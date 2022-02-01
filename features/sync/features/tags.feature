Feature: syncing the main branch syncs the tags

  Scenario: local tag gets pushed to the remote
    Given my repo has the following tags
      | NAME      | LOCATION |
      | local-tag | local    |
    And I am on the "main" branch
    When I run "git-town sync"
    Then my repo now has the following tags
      | NAME      | LOCATION      |
      | local-tag | local, remote |

  Scenario: tag on the remote branch gets pulled
    Given my repo has the following tags
      | NAME       | LOCATION |
      | remote-tag | remote   |
    And I am on the "main" branch
    When I run "git-town sync"
    Then my repo now has the following tags
      | NAME       | LOCATION      |
      | remote-tag | local, remote |

  Scenario: tag on a different branch gets pulled
    Given my repo has a remote tag "remote-tag" that is not on a branch
    And I am on the "main" branch
    When I run "git-town sync"
    Then my repo now has the following tags
      | NAME       | LOCATION      |
      | remote-tag | local, remote |
