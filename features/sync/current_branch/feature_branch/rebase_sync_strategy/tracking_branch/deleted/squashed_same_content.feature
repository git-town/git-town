Feature: delete a branch with deleted tracking branch whose changed files have the content of the parent
  (the local commits differ from the remote ones, so the branch is not in sync with its tracking branch)

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | shipped | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME | FILE CONTENT   |
      | shipped | local, origin | first version  | file      | first version  |
      | shipped | local, origin | second version | file      | second version |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "shipped"
    And I ran "git commit --amend --message 'second version rewritten'"
    And origin ships the "shipped" branch using the "squash-merge" ship-strategy
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | shipped | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git checkout shipped                              |
      | shipped | git checkout main                                 |
      | main    | git branch -D shipped                             |
    And no rebase is now in progress
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
