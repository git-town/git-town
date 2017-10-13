Feature: Syncing before creating the pull request

  As a developer
  I want my feature branch by synced before creating a pull request for it
  So that my reviewers see the most up-to-date version of my code and their review is accurate.


  Background:
    Given my repository has a feature branch named "parent-feature"
    And my repository has a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION | MESSAGE              | FILE NAME          |
      | main           | local    | local main commit    | local_main_file    |
      |                | remote   | remote main commit   | remote_main_file   |
      | parent-feature | local    | local parent commit  | local_parent_file  |
      |                | remote   | remote parent commit | remote_parent_file |
      | child-feature  | local    | local child commit   | local_child_file   |
      |                | remote   | remote child commit  | remote_child_file  |
    And I have "open" installed
    And my repo's remote origin is git@github.com:Originate/git-town.git
    And I am on the "child-feature" branch
    And my workspace has an uncommitted file
    When I run `git-town new-pull-request`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH         | COMMAND                                                                                    |
      | child-feature  | git fetch --prune                                                                          |
      |                | git add -A                                                                                 |
      |                | git stash                                                                                  |
      |                | git checkout main                                                                          |
      | main           | git rebase origin/main                                                                     |
      |                | git push                                                                                   |
      |                | git checkout parent-feature                                                                |
      | parent-feature | git merge --no-edit origin/parent-feature                                                  |
      |                | git merge --no-edit main                                                                   |
      |                | git push                                                                                   |
      |                | git checkout child-feature                                                                 |
      | child-feature  | git merge --no-edit origin/child-feature                                                   |
      |                | git merge --no-edit parent-feature                                                         |
      |                | git push                                                                                   |
      |                | git stash pop                                                                              |
      | <none>         | open https://github.com/Originate/git-town/compare/parent-feature...child-feature?expand=1 |
    And I see a new GitHub pull request for the "child-feature" branch against the "parent-feature" branch in the "Originate/git-town" repo in my browser
    And I am still on the "child-feature" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH         | LOCATION         | MESSAGE                                                                  | FILE NAME          |
      | main           | local and remote | remote main commit                                                       | remote_main_file   |
      |                |                  | local main commit                                                        | local_main_file    |
      | child-feature  | local and remote | local child commit                                                       | local_child_file   |
      |                |                  | remote child commit                                                      | remote_child_file  |
      |                |                  | Merge remote-tracking branch 'origin/child-feature' into child-feature   |                    |
      |                |                  | local parent commit                                                      | local_parent_file  |
      |                |                  | remote parent commit                                                     | remote_parent_file |
      |                |                  | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |                    |
      |                |                  | remote main commit                                                       | remote_main_file   |
      |                |                  | local main commit                                                        | local_main_file    |
      |                |                  | Merge branch 'main' into parent-feature                                  |                    |
      |                |                  | Merge branch 'parent-feature' into child-feature                         |                    |
      | parent-feature | local and remote | local parent commit                                                      | local_parent_file  |
      |                |                  | remote parent commit                                                     | remote_parent_file |
      |                |                  | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |                    |
      |                |                  | remote main commit                                                       | remote_main_file   |
      |                |                  | local main commit                                                        | local_main_file    |
      |                |                  | Merge branch 'main' into parent-feature                                  |                    |
