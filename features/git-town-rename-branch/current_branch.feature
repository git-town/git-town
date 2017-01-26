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
      | production | local and remote | prod commit |
      | feature    | local and remote | feat commit |


  Scenario: rename feature branch
    When I am on the "feature" branch
    And I run `git town-rename-branch renamed-feature`
    Then I end up on the "renamed-feature" branch
    And my repo is configured with perennial branches as "production"
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE     |
      | main            | local and remote | main commit |
      | production      | local and remote | prod commit |
      | renamed-feature | local and remote | feat commit |


  Scenario: rename perennial branch
    When I am on the "production" branch
    And I run `git town-rename-branch renamed-production -f`
    Then I end up on the "renamed-production" branch
    And my repo is configured with perennial branches as "renamed-production"
    And I have the following commits
      | BRANCH             | LOCATION         | MESSAGE     |
      | main               | local and remote | main commit |
      | feature            | local and remote | feat commit |
      | renamed-production | local and remote | prod commit |


  Scenario: undo rename branch
    Given I am on the "feature" branch
    And I run `git town-rename-branch renamed-feature`
    When I run `git town-rename-branch --undo`
    Then it runs the commands
        | BRANCH          | COMMAND                                     |
        | renamed-feature | git branch feature <%= sha 'feat commit' %> |
        |                 | git push -u origin feature                  |
        |                 | git push origin :renamed-feature            |
        |                 | git checkout feature                        |
        | feature         | git branch -d renamed-feature               |
    And I end up on the "feature" branch
    And my repo is configured with perennial branches as "production"
    And I have the following commits
      | BRANCH     | LOCATION         | MESSAGE     |
      | main       | local and remote | main commit |
      | feature    | local and remote | feat commit |
      | production | local and remote | prod commit |
