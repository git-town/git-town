Feature: Setting the pull strategy to "merge"

  As a developer
  I want to be able to configure Git Town to pull branches using `git merge`
  So that it behaves compatible with existing guidelines for my team.


  Scenario: Setting the pull strategy to "merge"
    When I run `git town config --pull-strategy merge`
    Then my repo is configured with the pull strategy "merge"
