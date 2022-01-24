Feature: git town-rename-branch: offline mode


  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE     |
      | main    | local, remote | main commit |
      | feature | local, remote | feat commit |
    And I am on the "feature" branch
    When I run "git-town rename-branch renamed-feature"


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                            |
      | feature         | git branch renamed-feature feature |
      |                 | git checkout renamed-feature       |
      | renamed-feature | git branch -D feature              |
    And I am now on the "renamed-feature" branch
    And my repo now has the following commits
      | BRANCH          | LOCATION      | MESSAGE     |
      | main            | local, remote | main commit |
      | feature         | remote        | feat commit |
      | renamed-feature | local         | feat commit |


  Scenario: undo rename branch
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH          | COMMAND                                    |
      | renamed-feature | git branch feature {{ sha 'feat commit' }} |
      |                 | git checkout feature                       |
      | feature         | git branch -D renamed-feature              |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE     |
      | main    | local, remote | main commit |
      | feature | local, remote | feat commit |
