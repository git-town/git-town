Feature: sync tags

  Scenario: local tag gets pushed to origin
    Given my repo has the tags
      | NAME      | LOCATION |
      | local-tag | local    |
    When I run "git-town sync"
    Then my repo now has the tags
      | NAME      | LOCATION      |
      | local-tag | local, origin |

  Scenario: tags on origin get pulled
    Given my repo has the tags
      | NAME       | LOCATION |
      | origin-tag | origin   |
    When I run "git-town sync"
    Then my repo now has the tags
      | NAME       | LOCATION      |
      | origin-tag | local, origin |

  Scenario: tag on a different branch gets pulled
    Given my repo has a remote tag "origin-tag" that is not on a branch
    When I run "git-town sync"
    Then my repo now has the tags
      | NAME       | LOCATION      |
      | origin-tag | local, origin |
