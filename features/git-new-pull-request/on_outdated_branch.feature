Feature: Syncing before creating the pull request

  As a developer
  I want that GT syncs my feature branch before creating a pull request for it
  So that my reviewers see the most up-to-date version of my code and their review is accurate.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           |
      | main    | local    | local main commit     | local_main_file     |
      |         | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      |         | remote   | remote feature commit | remote_feature_file |
    And I have "open" installed
    And my remote origin is git@github.com:Originate/git-town.git
    And I am on the "feature" branch
    When I run `git new-pull-request`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser
    And I am still on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME           |
      | main    | local and remote | remote main commit                                         | remote_main_file    |
      |         |                  | local main commit                                          | local_main_file     |
      | feature | local and remote | local feature commit                                       | local_feature_file  |
      |         |                  | remote feature commit                                      | remote_feature_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                     |
      |         |                  | remote main commit                                         | remote_main_file    |
      |         |                  | local main commit                                          | local_main_file     |
      |         |                  | Merge branch 'main' into feature                           |                     |

