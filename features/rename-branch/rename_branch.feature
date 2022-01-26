Feature: rename the current branch

  Background:
    Given my repo has a feature branch named "feature"
    And my repo has the perennial branch "production"
    And the following commits exist in my repo
      | BRANCH     | LOCATION      | MESSAGE     |
      | main       | local, remote | main commit |
      | feature    | local, remote | feat commit |
      | production | local, remote | prod commit |

  Scenario: rename feature branch
    Given I am on the "feature" branch
    When I run "git-town rename-branch renamed-feature"
    Then it runs the commands
      | BRANCH          | COMMAND                            |
      | feature         | git fetch --prune --tags           |
      |                 | git branch renamed-feature feature |
      |                 | git checkout renamed-feature       |
      | renamed-feature | git push -u origin renamed-feature |
      |                 | git push origin :feature           |
      |                 | git branch -D feature              |
    And I am now on the "renamed-feature" branch
    And the perennial branches are now configured as "production"
    And my repo now has the following commits
      | BRANCH          | LOCATION      | MESSAGE     |
      | main            | local, remote | main commit |
      | production      | local, remote | prod commit |
      | renamed-feature | local, remote | feat commit |

  Scenario: undo rename branch
    Given I am on the "feature" branch
    And I run "git-town rename-branch renamed-feature"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH          | COMMAND                                    |
      | renamed-feature | git branch feature {{ sha 'feat commit' }} |
      |                 | git push -u origin feature                 |
      |                 | git push origin :renamed-feature           |
      |                 | git checkout feature                       |
      | feature         | git branch -D renamed-feature              |
    And I am now on the "feature" branch
    And the perennial branches are now configured as "production"
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE     |
      | main       | local, remote | main commit |
      | feature    | local, remote | feat commit |
      | production | local, remote | prod commit |
