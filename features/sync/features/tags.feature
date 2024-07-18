Feature: sync tags

  Background:
    Given a Git repo clone

  Scenario: local tag gets pushed to origin
    Given the tags
      | NAME      | LOCATION |
      | local-tag | local    |
    When I run "git-town sync"
    Then these tags exist
      | NAME      | LOCATION      |
      | local-tag | local, origin |

  Scenario: tags on origin get pulled
    Given the tags
      | NAME       | LOCATION |
      | origin-tag | origin   |
    When I run "git-town sync"
    Then these tags exist
      | NAME       | LOCATION      |
      | origin-tag | local, origin |

  Scenario: tag on a different branch gets pulled
    Given a remote tag "origin-tag" not on a branch
    When I run "git-town sync"
    Then these tags exist
      | NAME       | LOCATION      |
      | origin-tag | local, origin |
