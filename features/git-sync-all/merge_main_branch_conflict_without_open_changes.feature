Feature: git-sync-all from the main branch

  Background:
    Given I have feature branches named "feature" and "local_feature"
    And my coworker has a feature branch named "remote_feature"
    And the following commits exist in my repository
      | branch         | location         | message               | file name           | file content           |
      | main           | local and remote | main commit           | conflicting_file    | main content           |
      | feature        | local and remote | feature commit        | conflicting_file    | feature content        |
      | local_feature  | local            | local feature commit  | local_feature_file  | local feature content  |
      | remote_feature | remote           | remote feature commit | remote_feature_file | remote feature content |
    And I am on the "main" branch
    When I run `git sync-all` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then I end up on the "feature" branch
    And my repo has a merge in progress
    And there are abort, skip, and continue scripts for "git sync-all"


  Scenario: aborting
    When I run `git sync-all --abort`
    Then I end up on the "main" branch
    And there are no abort, skip and continue scripts for "git sync-all" anymore
    And I have the following commits
      | branch         | location         | message               | files               |
      | main           | local and remote | main commit           | conflicting_file    |
      | feature        | local and remote | feature commit        | conflicting_file    |
      | local_feature  | local            | local feature commit  | local_feature_file  |
      | remote_feature | remote           | remote feature commit | remote_feature_file |


  Scenario: skipping
    When I run `git sync-all --skip`
    Then I end up on the "main" branch
    And there are no abort, skip and continue scripts for "git sync-all" anymore
    And I have the following commits
      | branch         | location         | message                                | files               |
      | main           | local and remote | main commit                            | conflicting_file    |
      | feature        | local and remote | feature commit                         | conflicting_file    |
      | local_feature  | local and remote | Merge branch 'main' into local_feature |                     |
      | local_feature  | local and remote | main commit                            | conflicting_file    |
      | local_feature  | local and remote | local feature commit                   | local_feature_file  |
      | remote_feature | remote           | remote feature commit                  | remote_feature_file |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync-all --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And my repo still has a merge in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync-all --continue`
    Then I end up on the "main" branch
    And there are no abort, skip and continue scripts for "git sync-all" anymore
    And I have the following commits
      | branch         | location         | message                                | files               |
      | main           | local and remote | main commit                            | conflicting_file    |
      | feature        | local and remote | Merge branch 'main' into feature       |                     |
      | feature        | local and remote | main commit                            | conflicting_file    |
      | feature        | local and remote | feature commit                         | conflicting_file    |
      | local_feature  | local and remote | Merge branch 'main' into local_feature |                     |
      | local_feature  | local and remote | main commit                            | conflicting_file    |
      | local_feature  | local and remote | local feature commit                   | local_feature_file  |
      | remote_feature | remote           | remote feature commit                  | remote_feature_file |


  Scenario: continuing after resolving conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit`
    And I run `git sync-all --continue`
    Then I end up on the "main" branch
    And there are no abort, skip and continue scripts for "git sync-all" anymore
    And I have the following commits
      | branch         | location         | message                                | files               |
      | main           | local and remote | main commit                            | conflicting_file    |
      | feature        | local and remote | Merge branch 'main' into feature       |                     |
      | feature        | local and remote | main commit                            | conflicting_file    |
      | feature        | local and remote | feature commit                         | conflicting_file    |
      | local_feature  | local and remote | Merge branch 'main' into local_feature |                     |
      | local_feature  | local and remote | main commit                            | conflicting_file    |
      | local_feature  | local and remote | local feature commit                   | local_feature_file  |
      | remote_feature | remote           | remote feature commit                  | remote_feature_file |
