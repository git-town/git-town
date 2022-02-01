Feature: syncing the main branch syncs the tags

  Scenario: Pushing tags
    Given my repo has the following tags
      | NAME      | LOCATION |
      | local-tag | local    |
    And I am on the "main" branch
    When I run "git-town sync"
    Then my repo now has the following tags
      | NAME      | LOCATION      |
      | local-tag | local, remote |

  Scenario: fetching tags on a pulled branch
    Given my repo has the following tags
      | NAME       | LOCATION |
      | remote-tag | remote   |
    And I am on the "main" branch
    When I run "git-town sync"
    Then my repo now has the following tags
      | NAME       | LOCATION      |
      | remote-tag | local, remote |

  Scenario: fetching tags not on a branch
    Given my repo has a remote tag "remote-tag" that is not on a branch
    And I am on the "main" branch
    When I run "git-town sync"
    Then my repo now has the following tags
      | NAME       | LOCATION      |
      | remote-tag | local, remote |
