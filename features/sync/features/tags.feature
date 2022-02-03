Feature: sync tags

  Scenario: local tag gets pushed to the remote
    Given my repo has the tags
      | NAME      | LOCATION |
      | local-tag | local    |
    When I run "git-town sync"
    Then my repo now has the tags
      | NAME      | LOCATION      |
      | local-tag | local, remote |

  Scenario: tag on the remote branch gets pulled
    Given my repo has the tags
      | NAME       | LOCATION |
      | remote-tag | remote   |
    When I run "git-town sync"
    Then my repo now has the tags
      | NAME       | LOCATION      |
      | remote-tag | local, remote |

  Scenario: tag on a different branch gets pulled
    Given my repo has a remote tag "remote-tag" that is not on a branch
    When I run "git-town sync"
    Then my repo now has the tags
      | NAME       | LOCATION      |
      | remote-tag | local, remote |
