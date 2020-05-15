Feature: git town-rename-branch: rename current branch implicitly

  As a developer wishing to rename the current branch
  I should be able reference the current branch implicitly
  So that I can perform my rename quickly.


  Background:
    Given my repository has a feature branch named "feature"
    And my repository has the perennial branch "production"
    And the following commits exist in my repository
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
    And I end up on the "renamed-feature" branch
    And the perennial branches are now configured as "production"
    And my repository now has the following commits
      | BRANCH          | LOCATION      | MESSAGE     |
      | main            | local, remote | main commit |
      | production      | local, remote | prod commit |
      | renamed-feature | local, remote | feat commit |


  Scenario: rename branch to itself
    Given I am on the "feature" branch
    When I run "git-town rename-branch feature"
    Then it runs no commands
    And it prints the error:
      """
      Cannot rename branch to current name.
      """
    And I end up on the "feature" branch
    And my repository is left with my original commits


  Scenario: rename perennial branch
    Given I am on the "production" branch
    When I run "git-town rename-branch renamed-production"
    Then it runs no commands
    And it prints the error:
      """
      "production" is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'.
      """


  Scenario: rename perennial branch (forced)
    Given I am on the "production" branch
    When I run "git-town rename-branch renamed-production --force"
    Then it runs the commands
      | BRANCH             | COMMAND                                  |
      | production         | git fetch --prune --tags                 |
      |                    | git branch renamed-production production |
      |                    | git checkout renamed-production          |
      | renamed-production | git push -u origin renamed-production    |
      |                    | git push origin :production              |
      |                    | git branch -D production                 |
    And I end up on the "renamed-production" branch
    And the perennial branches are now configured as "renamed-production"
    And my repository now has the following commits
      | BRANCH             | LOCATION      | MESSAGE     |
      | main               | local, remote | main commit |
      | feature            | local, remote | feat commit |
      | renamed-production | local, remote | prod commit |


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
    And I end up on the "feature" branch
    And the perennial branches are now configured as "production"
    And my repository now has the following commits
      | BRANCH     | LOCATION      | MESSAGE     |
      | main       | local, remote | main commit |
      | feature    | local, remote | feat commit |
      | production | local, remote | prod commit |
