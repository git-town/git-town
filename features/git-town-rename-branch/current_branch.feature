Feature: git town-rename-branch: rename current branch implicitly

  As a developer wishing to rename the current branch
  I should be able reference the current branch implicitly
  So that I can perform my rename quickly.


  Background:
    Given I have a feature branch named "feature"
    And I have a perennial branch named "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE     |
      | main       | local and remote | main commit |
      | feature    | local and remote | feat commit |
      | production | local and remote | prod commit |


  Scenario: rename feature branch
    Given I am on the "feature" branch
    When I run `gt rename-branch renamed-feature`
    Then it runs the commands
      | BRANCH          | COMMAND                            |
      | feature         | git fetch --prune                  |
      |                 | git branch renamed-feature feature |
      |                 | git checkout renamed-feature       |
      | renamed-feature | git push -u origin renamed-feature |
      |                 | git push origin :feature           |
      |                 | git branch -D feature              |
    And I end up on the "renamed-feature" branch
    And my repo is configured with perennial branches as "production"
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE     |
      | main            | local and remote | main commit |
      | production      | local and remote | prod commit |
      | renamed-feature | local and remote | feat commit |


  Scenario: rename branch to itself
    Given I am on the "feature" branch
    When I run `gt rename-branch feature`
    Then it runs no commands
    And I get the error "Cannot rename branch to current name."
    And I end up on the "feature" branch
    And I am left with my original commits


  Scenario: rename perennial branch
    Given I am on the "production" branch
    When I run `gt rename-branch renamed-production`
    Then it runs no commands
    And I get the error "production' is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'."


  Scenario: rename perennial branch (forced)
    Given I am on the "production" branch
    When I run `gt rename-branch renamed-production --force`
    Then it runs the commands
      | BRANCH             | COMMAND                                  |
      | production         | git fetch --prune                        |
      |                    | git branch renamed-production production |
      |                    | git checkout renamed-production          |
      | renamed-production | git push -u origin renamed-production    |
      |                    | git push origin :production              |
      |                    | git branch -D production                 |
    And I end up on the "renamed-production" branch
    And my repo is configured with perennial branches as "renamed-production"
    And I have the following commits
      | BRANCH             | LOCATION         | MESSAGE     |
      | main               | local and remote | main commit |
      | feature            | local and remote | feat commit |
      | renamed-production | local and remote | prod commit |


  Scenario: undo rename branch
    Given I am on the "feature" branch
    And I run `gt rename-branch renamed-feature`
    When I run `gt rename-branch --undo`
    Then it runs the commands
        | BRANCH          | COMMAND                                     |
        | renamed-feature | git branch feature <%= sha 'feat commit' %> |
        |                 | git push -u origin feature                  |
        |                 | git push origin :renamed-feature            |
        |                 | git checkout feature                        |
        | feature         | git branch -D renamed-feature               |
    And I end up on the "feature" branch
    And my repo is configured with perennial branches as "production"
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE     |
      | main       | local and remote | main commit |
      | feature    | local and remote | feat commit |
      | production | local and remote | prod commit |
